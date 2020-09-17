package v2api

import (
	"net/http"

	"github.com/sirupsen/logrus"
	stripeSDK "github.com/stripe/stripe-go"

	"gitlab.com/abios/user-svc/database/queries"
	"gitlab.com/abios/user-svc/errorspec"
	"gitlab.com/abios/user-svc/httphelper"
	"gitlab.com/abios/user-svc/logging"
	"gitlab.com/abios/user-svc/stripe"
	"gitlab.com/abios/user-svc/structs"
)

// gameIds is a helper function returns the keys of the specified game-map
func gameIds(games map[int64]bool) []int64 {
	keys := make([]int64, 0)
	for k := range games {
		keys = append(keys, k)
	}
	return keys
}

// GetSubscription is an HTTP handler for retrieving a customer's main/default
// subscription from Stripe
func GetSubscription(w http.ResponseWriter, r *http.Request) {
	reqLogger := logging.GetRequestLoggerFromContext(r.Context())
	customer, ok := r.Context().Value("customer").(structs.Customer)
	if !ok {
		httphelper.Fail(w, reqLogger,
			http.StatusInternalServerError,
			errorspec.MiddlewareFailure)
		return
	}
	// Fetch the list of products and games this customer is subscribed to
	productMapping, err := queries.GetCustomerProducts(customer.Id)
	if err != nil {
		reqLogger.WithError(err).Errorln("Unable to fetch customer's subscribed products")
		httphelper.Fail(w, reqLogger,
			http.StatusInternalServerError,
			errorspec.DatabaseFailure)
		return
	}
	// Create a mapping of Stripe product IDs to a set of game IDs. This mapping
	// shows what games are subscribed for each Stripe product, e.g. Series-
	// product for Dota & LOL and PbP-product for CSGO. Since a customer could
	// technically have two API clients with the same configuration, the Stripe ID
	// must map to a set of game IDs to avoid duplicate entries
	products := make(map[string]map[int64]bool)
	for _, p := range productMapping {
		if _, exists := products[p.StripeProductId]; !exists {
			products[p.StripeProductId] = make(map[int64]bool)
		}
		products[p.StripeProductId][p.GameId] = true
	}

	// Fetch customer information from Stripe
	stripeCustomer, err := stripe.GetCustomer(*customer.StripeId, nil)
	if err != nil {
		errspec := stripe.ParseError(err)
		if errspec.Critical {
			reqLogger.WithError(err).Errorln("Unable to get customer info via Stripe SDK")
		}
		httphelper.Fail(w, reqLogger,
			http.StatusInternalServerError,
			errspec.Code)
		return
	}

	// Create a subscription-item for each listed Stripe product
	subItems := make([]structs.SubscriptionItem, 0)
	for stripeId, games := range products {
		// Get the product's default plan from Stripe
		plan, err := stripe.GetDefaultPlan(stripeId)
		if err != nil {
			reqLogger.WithError(err).
				WithField("product_id", stripeId).
				Errorln("Unable to fetch default plan of Stripe Product")
			httphelper.Fail(w, reqLogger,
				http.StatusInternalServerError,
				errorspec.MissingDefaultPlan)
			return
		}
		// Currency consistency check
		if plan.Currency != stripeCustomer.Currency {
			reqLogger.WithError(err).
				WithFields(logrus.Fields{
					"customer_currency": string(stripeCustomer.Currency),
					"plan_currency":     string(plan.Currency),
				}).Errorln("Customer is subscribed to Stripe plans of invalid currencies")
			httphelper.Fail(w, reqLogger,
				http.StatusInternalServerError,
				errorspec.CurrencyMismatch)
			return
		}
		// Create a new subscription-item entry with the default-plan data
		subItems = append(subItems, structs.SubscriptionItem{
			Id:        stripeId,
			Name:      plan.Product.Name,
			UnitPrice: plan.Amount,
			Amount:    int64(len(games)) * plan.Amount,
			GameIds:   gameIds(games),
		})
	}

	// Attempt to query the status of this customer's default/main subscription on
	// Stripe. The status is set as "inactive" by default if there's no main sub
	var status string = "inactive"
	stripeSub, err := stripe.GetDefaultSubscription(*customer.StripeId)
	if err != nil {
		reqLogger.WithError(err).Errorln("Unable to find default Stripe subscription for customer")
		httphelper.Fail(w, reqLogger,
			http.StatusInternalServerError,
			errorspec.Unspecified)
		return
	}
	if stripeSub != nil {
		status = string(stripeSub.Status)
	}

	// Set the default discount to 0
	var discount int64 = 0
	if stripeCustomer.Discount != nil {
		discount = stripeCustomer.Discount.Coupon.AmountOff
	}

	// Calculate the final amount and the amount charged (incl. discount) and
	// return a Subscription struct with all the acquired information
	var amount int64
	for _, item := range subItems {
		amount += item.Amount
	}
	var charged int64 = amount - discount
	if charged < 0 {
		charged = 0
	}
	httphelper.Respond(w, reqLogger, http.StatusOK, structs.Subscription{
		Status:   status,
		Items:    subItems,
		Amount:   amount,
		Discount: discount,
		Charged:  charged,
		Currency: string(stripeCustomer.Currency),
	})
}

func ActivateSubscription(w http.ResponseWriter, r *http.Request) {
	reqLogger := logging.GetRequestLoggerFromContext(r.Context())
	customer, ok := r.Context().Value("customer").(structs.Customer)
	if !ok {
		httphelper.Fail(w, reqLogger,
			http.StatusInternalServerError,
			errorspec.MiddlewareFailure)
		return
	}
	// Attempt to find the default Stripe subscription for the customer
	stripeSub, err := stripe.GetDefaultSubscription(*customer.StripeId)
	if err != nil {
		reqLogger.WithError(err).Errorln("Unable to query default Stripe subscription for customer")
		httphelper.Fail(w, reqLogger,
			http.StatusInternalServerError,
			errorspec.Unspecified)
		return
	}
	// If there is already a default Stripe subscription, activation fails
	if stripeSub != nil {
		httphelper.Fail(w, reqLogger,
			http.StatusConflict,
			errorspec.SubscriptionAlreadyActive)
		return
	}
	// Fetch the list of products this customer is subscribed to
	products, err := queries.GetCustomerProducts(customer.Id)
	if err != nil {
		reqLogger.WithError(err).Errorln("Unable to fetch customer's subscribed products")
		httphelper.Fail(w, reqLogger,
			http.StatusInternalServerError,
			errorspec.DatabaseFailure)
		return
	}
	// Count the quantity of each product
	quantities := make(map[string]int64)
	for _, mapping := range products {
		quantities[mapping.StripeProductId] = quantities[mapping.StripeProductId] + 1
	}
	reqLogger.Debug(quantities)
	// Create a Stripe subscription-item for each product
	items := make([]*stripeSDK.SubscriptionItemsParams, 0)
	for productId := range quantities {
		plan, err := stripe.GetDefaultPlan(productId)
		if err != nil {
			reqLogger.WithError(err).
				WithField("product_id", productId).
				Errorln("Unable to fetch default plan of Stripe Product")
			httphelper.Fail(w, reqLogger,
				http.StatusInternalServerError,
				errorspec.MissingDefaultPlan)
			return
		}
		items = append(items, &stripeSDK.SubscriptionItemsParams{
			Plan:     &plan.ID,
			Quantity: &quantities[productId],
		})
	}
	// Create the default-subscription
	_, err = stripe.AddSubscription(&stripeSDK.SubscriptionParams{
		Customer: customer.StripeId,
		Items:    items,
	})
	if err != nil {
		errspec := stripe.ParseError(err)
		if errspec.Critical {
			reqLogger.WithError(err).Errorln("Unable to create Stripe subscription")
		}
		httphelper.Fail(w, reqLogger,
			http.StatusInternalServerError,
			errspec.Code)
		return
	}

	httphelper.Respond(w, reqLogger, http.StatusNoContent)
}

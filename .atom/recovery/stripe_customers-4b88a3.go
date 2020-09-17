package v2api

import (
	"net/http"

	"github.com/sirupsen/logrus"
	stripePkg "github.com/stripe/stripe-go"

	"gitlab.com/abios/user-svc/database"
	"gitlab.com/abios/user-svc/database/queries"
	"gitlab.com/abios/user-svc/errorspec"
	"gitlab.com/abios/user-svc/httphelper"
	"gitlab.com/abios/user-svc/logging"
	"gitlab.com/abios/user-svc/stripe"
)

// createCustomer is used for creating a new customer in stripe
func createCustomer(reqLogger *logrus.Entry, id int64, name string) (*stripePkg.Customer, error) {
	params := stripePkg.CustomerParams{
		Name: &name,
	}
	// Create new stripe customer
	customerStripe, err := stripe.CreateCustomer(&params)
	if err != nil {
		reqLogger.WithError(err).Errorln("Couldn't create new stripe customer")
		return customerStripe, err
	}
	// Insert new stripe ID for customer
	q := "UPDATE Customer SET stripe_id=? WHERE id=?"
	_, err = database.OAuthDB.Exec(q, customerStripe.ID, id)
	if err != nil {
		reqLogger.WithError(err).Errorln("Couldn't insert new stripe ID for customer")
		// If update fails, remove stripe customer and respond back with an error
		_, err = stripe.DeleteCustomer(customerStripe.ID)
		if err != nil {
			reqLogger.WithError(err).
				WithField("stripe-id", customerStripe.ID).
				Errorln("Manual delete is required for stripe customer")
		}
		return customerStripe, err
	}
	return customerStripe, nil
}

func GetStripeCustomer(w http.ResponseWriter, r *http.Request) {
	reqLogger := logging.GetRequestLoggerFromContext(r.Context())
	customer, ok := r.Context().Value("customer").(structs.Customer)
	if !ok {
		httphelper.Fail(w, reqLogger,
			http.StatusInternalServerError,
			errorspec.MiddlewareFailure)
		return
	}

	// Get customer info from stipe
	stripeCustomer, err := stripe.GetCustomer(customer.StripeId, nil)
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

	httphelper.Respond(w, reqLogger, http.StatusOK, stripeCustomer)
}

func UpdateStripeCustomer(w http.ResponseWriter, r *http.Request) {
	reqLogger := logging.GetRequestLoggerFromContext(r.Context())

	// Extract customer ID
	customerId, ok := r.Context().Value("customer_id").(int64)
	if !ok {
		httphelper.Fail(w, reqLogger,
			http.StatusInternalServerError,
			errorspec.MiddlewareFailure)
		return
	}

	// Get customer data
	customer, err := queries.GetCustomer(customerId)
	if err != nil {
		reqLogger.WithError(err).Errorln("Couldn't get customer information from database")
		httphelper.Fail(w, reqLogger,
			http.StatusInternalServerError,
			errorspec.DatabaseFailure)
		return
	}

	// Check if customer already has a stripe ID
	if customer.StripeId != nil {
		// Stripe customer already exists
		httphelper.Fail(w, reqLogger, http.StatusConflict)
		return
	}

	// Create new stripe customer
	customerStripe, err := createCustomer(reqLogger, customerId, customer.Name)
	if err != nil {
		errspec := stripe.ParseError(err)
		if errspec.Critical {
			reqLogger.WithError(err).Errorln("Unable to create a new stripe customer via Stripe SDK")
		}
		httphelper.Fail(w, reqLogger,
			http.StatusInternalServerError,
			errspec.Code)
		return
	}

	httphelper.Respond(w, reqLogger, http.StatusOK, customerStripe)
}

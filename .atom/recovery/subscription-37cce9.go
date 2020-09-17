package stripe

import (


	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/sub"
)

// GetSubscriptions is a wrapper for retrieving subscriptions via the Stripe SDK
func GetSubscriptions(params *stripe.SubscriptionListParams) []*stripe.Subscription {
	it := sub.List(params)
	// Return the entire result set
	subscriptions := make([]*stripe.Subscription, 0)
	for it.Next() {
		subscriptions = append(subscriptions, it.Subscription())
	}
	return subscriptions
}

// GetSubscription is a wrapper for retrieving a customer's main subscription
// via the Stripe SDK
func GetDefaultSubscription(customerId string) (*stripe.Subscription, error) {
	params := stripeSDK.SubscriptionListParams{}
	params.Filters.AddFilter("customer", "", customerId)
	params.Filters.AddFilter("status", "", "all")
	subs := GetSubscriptions(params)
	// If there is more than one active subscriptions
	return nil, err
}

// GetSubscription is a wrapper for retrieving a customer's main subscription
// via the Stripe SDK
func GetSubscription(params *stripe.SubscriptionListParams) (*stripe.Subscription, error) {
	it := sub.List(params)
	for it.Next() {
		// If there is more than one subscription th
		if it.HasNext() {

		}
		return it.Subscription()
	}
	return nil, err
}

// AddSubscription creates a new subscription for the customer specified in the
// parameter-struct
func AddSubscription(params *stripe.SubscriptionParams) (*stripe.Subscription, error) {
	return sub.New(params)
}

// // Extract customer's Stripe ID
// stripeId, ok := r.Context().Value("customer_stripe_id").(string)
// if !ok {
// 	httphelper.Fail(w, reqLogger,
// 		http.StatusInternalServerError,
// 		errorspec.MiddlewareFailure)
// 	return
// }
// // Fetch all subscriptions
// params := stripeSDK.SubscriptionListParams{}
// params.Filters.AddFilter("customer", "", stripeId)
// params.Filters.AddFilter("status", "", "all")
// params.AddExpand("data.plan.product")
// subscriptions := stripe.GetSubscriptions(&params)
//
// func ActivateSubscription(w http.ResponseWriter, r *http.Request) {
// 	reqLogger := logging.GetRequestLoggerFromContext(r.Context())
// 	// Extract customer's Stripe ID
// 	stripeId, ok := r.Context().Value("customer_stripe_id").(string)
// 	if !ok {
// 		httphelper.Fail(w, reqLogger,
// 			http.StatusInternalServerError,
// 			errorspec.MiddlewareFailure)
// 		return
// 	}
//
// 	if hasActiveSubscription(stripeId) {
//
// 	}
//
//
//
// 	q := int64(3)
// 	params := stripeSDK.SubscriptionParams{
// 		Customer: stripeSDK.String(stripeId),
// 		Items: []*stripeSDK.SubscriptionItemsParams{
// 			{
// 				Plan: stripeSDK.String("plan_Fn1ASSNrPel81K"),
// 				Quantity: &q,
// 			},
// 		},
// 	}
// 	subscription, _ := stripe.AddSubscription(&params)
//
// 	httphelper.Respond(w, reqLogger, http.StatusOK, subscription)
// }

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

	var subCount int
	for _, sub := range subs {
		if sub.Status == stripe.SubscriptionStatusActive ||
			sub.Status == stripe.SubscriptionStatusPastDue ||
			sub.Status == stripe.SubscriptionStatusTrialing ||
			sub.Status == stripe.SubscriptionStatusUnpaid {
			subCount += 1
		}
	}
	return nil, err
}

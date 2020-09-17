package stripe

import (
	"errors"

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
	params := stripe.SubscriptionListParams{}
	params.Filters.AddFilter("customer", "", customerId)
	params.Filters.AddFilter("status", "", "all")
	subs := GetSubscriptions(&params)

	// Helper function for determining if a subscription is active. In this
	// interpretation a subscription is also considered active if it is unpaid,
	// past due date or undergoing a trial period
	isActive := func(sub *stripe.Subscription) bool {
		return sub.Status == stripe.SubscriptionStatusUnpaid ||
			sub.Status == stripe.SubscriptionStatusTrialing ||
			sub.Status == stripe.SubscriptionStatusPastDue ||
			sub.Status == stripe.SubscriptionStatusActive
	}
	// Attempt to find the single active subscription
	var mainSub *stripe.Subscription
	for _, sub := range subs {
		if isActive(sub) {
			// If there is more than one active subscription the request becomes
			// ambiguous and an error is raised
			if mainSub != nil {
				return nil, errors.New("Customer has more than one active subscription")
			} else {
				mainSub = sub
			}
		}
	}
	return mainSub, nil
}

func AddSubscription(params stripe.SubscriptionParams) {

}

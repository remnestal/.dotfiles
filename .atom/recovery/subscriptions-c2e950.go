package handlers

import (
	"net/http"

	"gitlab.com/abios/api-dash-backend/logging"
	"gitlab.com/abios/api-dash-backend/structs"
	"gitlab.com/abios/api-dash-backend/structs/exported"
	"gitlab.com/abios/api-dash-backend/usersvc"
)

// GetSubscription is an HTTP handler for retrieving information about the
// customer's subscription
func GetSubscription(w http.ResponseWriter, r *http.Request) {
	reqLogger := logging.GetRequestLoggerFromContext(r.Context())

	// Infer customer ID from context
	user, ok := r.Context().Value("User").(structs.UserInternal)
	if !ok {
		reqLogger.Errorln("Unable to get internal user")
		WriteError(w, reqLogger, http.StatusInternalServerError)
		return
	}
	// Fetch subscription via user-svc
	sub, err := usersvc.GetSubscription(user.Customer.Id)
	if err != nil {
		// No kind of error type is expected to be returned from the user-svc, thus
		// every kind of error is logged and listed as an internal server error
		reqLogger.WithError(err).Errorln("Unable to get subscription via user-svc")
		WriteError(w, reqLogger,
			http.StatusInternalServerError,
			"Unable to fetch payment-plan info! Please try again later.")
		return
	}

	// Create a dashboard version of all subscription-items from the user-svc
	items := make([]exported.SubscriptionItem, 0)
	for _, item := range sub.Items {
		items = append(items, exported.SubscriptionItem{
			Name:      item.Name,
			Id:        item.Id,
			UnitPrice: item.UnitPrice,
			Amount:    item.Amount,
			GameIds:   item.GameIds,
		})
	}

	// Return a ported version of the user-svc Subscription struct
	Respond(w, reqLogger, http.StatusOK, exported.Subscription{
		Status:   sub.Status,
		Discount: sub.Discount,
		Amount:   sub.Amount,
		Charged:  sub.Charged,
		Currency: sub.Currency,
		Items:    items,
	})
}

// StartSubscription activates the subscription of the customer associated to
// the user invoking this request. If there is already an active subscription,
// an error is returned instead
func StartSubscription(w http.ResponseWriter, r *http.Request) {
	reqLogger := logging.GetRequestLoggerFromContext(r.Context())

	// Infer customer ID from context
	user, ok := r.Context().Value("User").(structs.UserInternal)
	if !ok {
		reqLogger.Errorln("Unable to get internal user")
		WriteError(w, reqLogger, http.StatusInternalServerError)
		return
	}

	// Fetch subscription via user-svc
	if err := usersvc.StartSubscription(user.Customer.Id); err != nil {
		// Write an appropriate error message if an anticipated error type is found
		switch err {
		case errorspec.SubscriptionAlreadyActive:
			WriteError(w, reqLogger,
				http.StatusConflict,
				"Your subscription is already active")
		default:
			reqLogger.WithError(err).Errorln("An unknown error occured in the user-svc")
			WriteError(w, reqLogger,
				http.StatusInternalServerError,
				"Unable to start your subscription, please try again later!")
		}
		return
	}

	Respond(w, reqLogger, http.StatusNoContent)
}

func StopSubscription(w http.ResponseWriter, r *http.Request) {
	reqLogger := logging.GetRequestLoggerFromContext(r.Context())

	// Infer customer ID from context
	user, ok := r.Context().Value("User").(structs.UserInternal)
	if !ok {
		reqLogger.Errorln("Unable to get internal user")
		WriteError(w, reqLogger, http.StatusInternalServerError)
		return
	}

	// Fetch subscription via user-svc
	if err := usersvc.StopSubscription(user.Customer.Id); err != nil {
		// Write an appropriate error message if an anticipated error type is found
		switch err {
		case errorspec.SubscriptionAlreadyActive:
			WriteError(w, reqLogger,
				http.StatusConflict,
				"Your subscription is already active")
		default:
			reqLogger.WithError(err).Errorln("An unknown error occured in the user-svc")
			WriteError(w, reqLogger,
				http.StatusInternalServerError,
				"Unable to start your subscription, please try again later!")
		}
		return
	}

	Respond(w, reqLogger, http.StatusNoContent)
}

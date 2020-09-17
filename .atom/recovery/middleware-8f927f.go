package main

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"gitlab.com/abios/user-svc/database/queries"
	"gitlab.com/abios/user-svc/logging"
	"gitlab.com/abios/user-svc/report"
	"gitlab.com/abios/user-svc/stripe"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Set the default response content-type to JSON
func JSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// Adds a logger to the request context.
func SetupRequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ajax_endpoint := vars["endpoint"]

		reqLogger := svcLogger.WithFields(logrus.Fields{
			"endpoint": ajax_endpoint,
		})

		c := logging.AddRequestLoggerToContext(r.Context(), reqLogger)

		next.ServeHTTP(w, r.WithContext(c))
	})
}

// CustomerIdParser is an HTTP middlware handler that parses an expected URL
// parameter `customer_id` as int64 and sets it in the context for the next HTTP
// handler in the pipeline. If the argument can't be parsed as an int64 an
// MalformedCustomerId error is returned
func CustomerIdParser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqLogger := logging.GetRequestLoggerFromContext(r.Context())
		customerId, err := strconv.ParseInt(mux.Vars(r)["customer_id"], 10, 64)
		if err != nil {
			report.Write(w, reqLogger,
				http.StatusBadRequest,
				report.ErrorSpec{MalformedCustomerId: true})
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "customer_id", customerId)
		ctx = logging.AddRequestLoggerToContext(ctx, reqLogger.WithField("customer_id", customerId))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CustomerStripeId is an HTTP middleware handler that fetches the Stripe ID of
// the specified customer and attaches it to the request-context. If the
// customer can't be fetched from database, an appropriate error is returned
func CustomerStripeId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqLogger := logging.GetRequestLoggerFromContext(r.Context())
		// Extract customer ID
		customerId, ok := r.Context().Value("customer_id").(int64)
		if !ok {
			reqLogger.
				WithField("delegated-by", "CustomerStripeId middleware handler").
				Errorln("Customer ID can't be pulled from context")
			report.Write(w, reqLogger,
				http.StatusInternalServerError,
				report.ErrorSpec{MiddlewareFailure: true})
			return
		}
		// Fetch the customer from the database
		customer, err := queries.GetCustomer(customerId)
		if err != nil {
			if err == sql.ErrNoRows {
				report.Write(w, reqLogger,
					http.StatusNotFound,
					report.ErrorSpec{CustomerNotFound: true})
				return
			} else {
				reqLogger.WithError(err).
					WithField("customer-id", customerId).
					Errorln("Unable to get customer from database")
				report.Write(w, reqLogger,
					http.StatusInternalServerError,
					report.ErrorSpec{Unknown: true})
				return
			}
		}
		// Verify that the customer has a Stripe ID
		if customer.StripeId == nil {
			report.Write(w, reqLogger,
				http.StatusNotFound,
				report.ErrorSpec{MissingStripeId: true})
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "customer_stripe_id", *customer.StripeId)
		ctx = logging.AddRequestLoggerToContext(ctx, reqLogger.WithField("customer_stripe_id", *customer.StripeId))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ValidateCardOwner is an HTTP middleware handler which verifies that the
// customer specified in the request is in fact the owner of the specified card.
// This check is necessary to ensure that e.g. a customer can't issue removal of
// another customer's card
func VerifyCardOwner(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqLogger := logging.GetRequestLoggerFromContext(r.Context())
		// Extract customer's Stripe ID
		stripeId, ok := r.Context().Value("customer_stripe_id").(string)
		if !ok {
			reqLogger.
				WithField("delegated-by", "ValidateCardOwner middleware handler").
				Errorln("Customer's Stripe ID can't be pulled from context")
			report.Write(w, reqLogger,
				http.StatusInternalServerError,
				report.ErrorSpec{MiddlewareFailure: true})
			return
		}
		// Verify that the specified customer is the owner of the card
		paymentMethodID := mux.Vars(r)["card_id"]
		pm, err := stripe.GetPaymentMethod(paymentMethodID, nil)
		if err != nil {
			reqLogger.WithError(err).
				WithField("payment-method-id", paymentMethodID).
				Errorln("Unable to find specified payment method")
			report.Write(w, reqLogger,
				http.StatusInternalServerError,
				report.ErrorSpec{Unknown: true})
			return
		}
		// If there is no customer associated with the payment method, it has
		// aldready been detached.
		// FIXME: Should this be logged? Or is it "expected behavior"
		if pm.Customer == nil {
			report.Write(w, reqLogger,
				http.StatusInternalServerError,
				report.ErrorSpec{Unknown: true})
			return
		}
		// Check if the specified customer is the owner of the payment method
		if pm.Customer.ID != stripeId {
			report.Write(w, reqLogger,
				http.StatusForbidden,
				report.ErrorSpec{CustomerNotPermitted: true})
			return
		}
		next.ServeHTTP(w, r)
	})
}

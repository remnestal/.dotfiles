package v2api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	stripeSDK "github.com/stripe/stripe-go"

	"gitlab.com/abios/user-svc/database/queries"
	"gitlab.com/abios/user-svc/errorspec"
	"gitlab.com/abios/user-svc/httphelper"
	"gitlab.com/abios/user-svc/logging"
	"gitlab.com/abios/user-svc/stripe"
	"gitlab.com/abios/user-svc/structs"
)

// ListCards is an HTTP handler that returns a list of the specified customer's
// debit/credit cards listed in Stripe
func ListCards(w http.ResponseWriter, r *http.Request) {
	reqLogger := logging.GetRequestLoggerFromContext(r.Context())
	customer, ok := r.Context().Value("customer").(structs.Customer)
	if !ok {
		httphelper.Fail(w, reqLogger,
			http.StatusInternalServerError,
			errorspec.MiddlewareFailure)
		return
	}

	// Get all payment methods from Stripe
	params := stripeSDK.PaymentMethodListParams{
		Customer: stripeSDK.String(*customer.StripeId),
		Type:     stripeSDK.String(string(stripeSDK.PaymentMethodTypeCard)),
	}
	params.AddExpand("data.customer")
	paymentMethods := stripe.GetPaymentMethods(&params)

	// Prepare the list of cards to return
	cards := make([]structs.PaymentMethod, 0)
	for _, pm := range paymentMethods {
		// This check should not be necessary unless Stripe's API is broken, but
		// checking anyway to avoid panic
		if pm.Card == nil {
			reqLogger.
				WithField("paymentmethod-id", pm.ID).
				Errorln("No card associated with payment method")
			httphelper.Fail(w, reqLogger, http.StatusInternalServerError)
			return
		}
		cards = append(cards, structs.PaymentMethod{
			Id:      pm.ID,
			Default: stripe.IsDefaultPaymentMethod(*pm),
			CardholderName: pm.BillingDetails.Name,
			Card: structs.Card{
				Fingerprint: pm.Card.Fingerprint,
				ExpMonth:    pm.Card.ExpMonth,
				ExpYear:     pm.Card.ExpYear,
				Last4:       pm.Card.Last4,
				Brand:       string(pm.Card.Brand),
			},
		})
	}

	httphelper.Respond(w, reqLogger, http.StatusOK, cards)
}

// AddCard is an HTTP handler that creates a credit/debit card in Stripe and
// attaches it to the specified customer; as long as all parameters are sound
// and that the customer exists with a valid Stripe ID
func AddCard(w http.ResponseWriter, r *http.Request) {
	reqLogger := logging.GetRequestLoggerFromContext(r.Context())
	// Decode POST body
	var params structs.CardParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		httphelper.Fail(w, reqLogger,
			http.StatusBadRequest,
			errorspec.MalformedRequestBody)
		return
	}
	// Extract customer struct
	customer, ok := r.Context().Value("customer").(structs.Customer)
	if !ok {
		httphelper.Fail(w, reqLogger,
			http.StatusInternalServerError,
			errorspec.MiddlewareFailure)
		return
	}

	// Helper function for converting int64 to *string
	convert := func(i int64) *string {
		s := strconv.FormatInt(i, 10)
		return &s
	}
	// Create the card resource
	unattached, err := stripe.NewCard(&stripeSDK.PaymentMethodParams{
		Type: stripeSDK.String(string(stripeSDK.PaymentMethodTypeCard)),
		Card: &stripeSDK.PaymentMethodCardParams{
			CVC:      convert(params.CVC),
			Number:   convert(params.Number),
			ExpYear:  convert(params.ExpYear),
			ExpMonth: convert(params.ExpMonth),
		},
		BillingDetails: &stripeSDK.BillingDetailsParams{
			Name: params.CardholderName,
		},
	})
	if err != nil {
		errspec := stripe.ParseError(err)
		if errspec.Critical {
			reqLogger.WithError(err).Errorln("Unable to create new card in stripe")
		}
		httphelper.Fail(w, reqLogger,
			http.StatusInternalServerError,
			errspec.Code)
		return
	}

	// Check if the a card with the same fingerprint (same card number) is already
	// registered to a customer in the database
	if card, err := queries.GetCardByFingerprint(unattached.Card.Fingerprint); err != nil {
		reqLogger.WithError(err).
			WithField("card-fingerprint", unattached.Card.Fingerprint).
			Errorln("Unable to fetch card mapping from the database")
		httphelper.Fail(w, reqLogger,
			http.StatusInternalServerError,
			errorspec.DatabaseFailure)
		return
	} else if card != nil {
		// A card with that fingerprint already exists
		httphelper.Fail(w, reqLogger,
			http.StatusConflict,
			errorspec.DuplicateCard)
		return
	}

	// Attach the card to the Stripe customer resource
	attached, err := stripe.AttachCard(unattached.ID, &stripeSDK.PaymentMethodAttachParams{
		Customer: customer.StripeId,
	})
	if err != nil {
		errspec := stripe.ParseError(err)
		if errspec.Critical {
			reqLogger.WithError(err).WithFields(logrus.Fields{
				"payment-method-id": unattached.ID,
				"card-fingerprint":  unattached.Card.Fingerprint,
			}).Errorln("Unable to attach Stripe card resource")
		}
		httphelper.Fail(w, reqLogger,
			http.StatusInternalServerError,
			errspec.Code)
		return
	}

	// Add the card-mapping to the database if the card was successfully added
	// attached to the Customer in Stripe
	mapping := structs.CardMapping{
		CustomerId:  customer.Id,
		Fingerprint: attached.Card.Fingerprint,
	}
	if err := queries.AddCard(mapping); err != nil {
		reqLogger.WithError(err).WithFields(logrus.Fields{
			"payment-method-id": attached.ID,
			"card-fingerprint":  mapping.Fingerprint,
		}).Errorln("Unable to add card fingerprint mapping to database")
		httphelper.Fail(w, reqLogger,
			http.StatusInternalServerError,
			errorspec.DatabaseFailure)
		return
	}

	httphelper.Respond(w, reqLogger, http.StatusOK, structs.PaymentMethod{
		Id:      attached.ID,
		Default: false, // New cards are never set as default
		CardholderName: attached.BillingDetails.Name,
		Card: structs.Card{
			Fingerprint: attached.Card.Fingerprint,
			ExpMonth:    attached.Card.ExpMonth,
			ExpYear:     attached.Card.ExpYear,
			Last4:       attached.Card.Last4,
			Brand:       string(attached.Card.Brand),
		},
	})
}

// DeleteCard is an HTTP handler that removes the specified card from use in
// Stripe. The card-customer mapping is however only soft-deleted in the
// database.
func DeleteCard(w http.ResponseWriter, r *http.Request) {
	reqLogger := logging.GetRequestLoggerFromContext(r.Context())
	customer, ok := r.Context().Value("customer").(structs.Customer)
	if !ok {
		httphelper.Fail(w, reqLogger,
			http.StatusInternalServerError,
			errorspec.MiddlewareFailure)
		return
	}
	// Detach the payment method from its customer
	paymentMethodID := mux.Vars(r)["card_id"]
	pm, err := stripe.DetachPaymentMethod(paymentMethodID, nil)
	if err != nil {
		errspec := stripe.ParseError(err)
		if errspec.Critical {
			reqLogger.WithError(err).
				WithField("payment-method-id", paymentMethodID).
				Errorln("Unable to detach payment method")
		}
		httphelper.Fail(w, reqLogger,
			http.StatusInternalServerError,
			errspec.Code)
		return
	}
	// Soft-delete the card mapping in the database
	mapping := structs.CardMapping{
		CustomerId:  customer.Id,
		Fingerprint: pm.Card.Fingerprint,
	}
	if err := queries.DeleteCard(mapping); err != nil {
		reqLogger.WithError(err).WithFields(logrus.Fields{
			"payment-method-id": pm.ID,
			"card-fingerprint":  pm.Card.Fingerprint,
		}).Errorln("Unable to soft-delete card fingerprint mapping in database")
		httphelper.Fail(w, reqLogger,
			http.StatusInternalServerError,
			errorspec.DatabaseFailure)
		return
	}

	httphelper.Respond(w, reqLogger, http.StatusOK)
}

// SetDefaultCard is an HTTP handler that attempts to set the card (payment
// method) specified by URL parameter {card_id} as the default payment method
func SetDefaultCard(w http.ResponseWriter, r *http.Request) {
	reqLogger := logging.GetRequestLoggerFromContext(r.Context())
	customer, ok := r.Context().Value("customer").(structs.Customer)
	if !ok {
		httphelper.Fail(w, reqLogger,
			http.StatusInternalServerError,
			errorspec.MiddlewareFailure)
		return
	}

	// Update customer with the specified payment-method as default
	paymentMethodID := mux.Vars(r)["card_id"]
	params := stripeSDK.CustomerParams{
		InvoiceSettings: &stripeSDK.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: &paymentMethodID,
		},
	}
	_, err := stripe.UpdateCustomer(*customer.StripeId, &params)
	if err != nil {
		errspec := stripe.ParseError(err)
		if errspec.Critical {
			reqLogger.WithError(err).
				WithField("payment-method-id", paymentMethodID).
				Errorln("Unable to update Stripe customer with new default payment method")
		}
		httphelper.Fail(w, reqLogger,
			http.StatusInternalServerError,
			errspec.Code)
		return
	}

	httphelper.Respond(w, reqLogger, http.StatusOK)
}

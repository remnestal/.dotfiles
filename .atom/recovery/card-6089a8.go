package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"gitlab.com/abios/api-dash-backend/errorspec"
	"gitlab.com/abios/api-dash-backend/logging"
	"gitlab.com/abios/api-dash-backend/structs"
	"gitlab.com/abios/api-dash-backend/structs/exported"
	"gitlab.com/abios/api-dash-backend/usersvc"
	usersvcStructs "gitlab.com/abios/user-svc/structs"
)

// ListCards is an HTTP handler for retrieving the debit/credit cards of a
// customer, as listed in Stripe, via the user-svc
func ListCards(w http.ResponseWriter, r *http.Request) {
	reqLogger := logging.GetRequestLoggerFromContext(r.Context())

	// Infer customer ID from context
	user, ok := r.Context().Value("User").(structs.UserInternal)
	if !ok {
		reqLogger.Errorln("Unable to get internal user")
		WriteError(w, reqLogger, http.StatusInternalServerError)
		return
	}

	// Fetch cards via user-svc
	paymentMethods, err := usersvc.GetCards(user.Customer.Id)
	if err != nil {
		// No kind of error type is expected to be returned from the user-svc, thus
		// every kind of error is logged and listed as an internal server error
		reqLogger.WithError(err).Errorln("Unable to get cards via user-svc")
		WriteError(w, reqLogger, http.StatusInternalServerError)
		return
	}

	// Return a list of safe-to-share Card represntations if the request to the
	// user-svc was successful and a list of cards was returned
	expCards := make([]exported.Card, 0)
	for _, pm := range paymentMethods {
		expCards = append(expCards, exported.Card{
			Id:       pm.Id,
			Brand:    pm.Card.Brand,
			ExpMonth: pm.Card.ExpMonth,
			ExpYear:  pm.Card.ExpYear,
			LastFour: pm.Card.Last4,
			Default:  pm.Default,
			CardholderName:  pm.CardholderName,
		})
	}
	Respond(w, reqLogger, http.StatusOK, expCards)
}

// AddCard is an HTTP handler that creates a payment-card in Stripe, via the
// user-svc. If the creation-request succeded, information about the card is
// returned. Otherwise a JSON object with an error message is sent
func AddCard(w http.ResponseWriter, r *http.Request) {
	reqLogger := logging.GetRequestLoggerFromContext(r.Context())

	// Parse input parameters before propagating request to user-svc
	var input usersvcStructs.CardParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		WriteError(w, reqLogger,
			http.StatusBadRequest,
			"Malformed input data")
		return
	}

	// Infer customer ID from context
	user, ok := r.Context().Value("User").(structs.UserInternal)
	if !ok {
		reqLogger.Errorln("Unable to get internal user")
		WriteError(w, reqLogger, http.StatusInternalServerError)
		return
	}

	// Fetch billing information from user-svc
	billingData, err := usersvc.GetBilling(user.Customer.Id)
	if err != nil {
		reqLogger.WithError(err).Errorln("Unable to get customer billing info")
		WriteError(w, reqLogger, http.StatusInternalServerError)
		return
	}
	var billing usersvcStructs.StripeBilling
	if err = json.Unmarshal(billingData, &billing); err != nil {
		reqLogger.WithError(err).
			WithField("body", string(billingData)).
			Errorln("Unable to parse billing info")
		WriteError(w, reqLogger, http.StatusInternalServerError)
		return
	}

	// Check that billing information is OK
	// FIXME: this check should probably be more rigorous?
	if billing.Line1 == nil || billing.Email == "" {
		WriteError(w, reqLogger,
			http.StatusBadRequest,
			"Insufficient billing information")
		return
	}

	// Turn the card parameter-struct into bytes and forward it to the user-svc
	bytes, err := json.Marshal(input)
	if err != nil {
		reqLogger.WithError(err).Errorln("Unable to marshal card data")
		WriteError(w, reqLogger, http.StatusInternalServerError)
		return
	}
	paymentMethod, err := usersvc.CreateCard(user.Customer.Id, bytes)
	if err != nil {
		// Write an appropriate error message if an anticipated error type is found
		switch err {
		case errorspec.DuplicateCard:
			WriteError(w, reqLogger, http.StatusConflict, "Card is already in use")
		case errorspec.CardExpired:
			WriteError(w, reqLogger, http.StatusBadRequest, "Card has already expired")
		case errorspec.InvalidCardCVC:
			WriteError(w, reqLogger, http.StatusBadRequest, "Card's CVC value is invalid")
		case errorspec.InvalidCardNumber:
			WriteError(w, reqLogger, http.StatusBadRequest, "Card's number is invalid")
		case errorspec.InvalidCardExpMonth:
			WriteError(w, reqLogger, http.StatusBadRequest, "Card's expiration month is invalid")
		case errorspec.InvalidCardExpYear:
			WriteError(w, reqLogger, http.StatusBadRequest, "Card's expiration year is invalid")
		default:
			reqLogger.WithError(err).Errorln("An unknown error occured in the user-svc")
			WriteError(w, reqLogger, http.StatusInternalServerError)
		}
		return
	}

	// Get the list of cards to determine if the new card should be made default
	cards, err := usersvc.GetCards(user.Customer.Id)
	if err != nil {
		reqLogger.WithError(err).Errorln("Unable to get cards via user-svc")
		WriteError(w, reqLogger, http.StatusInternalServerError)
		return
	}
	// If the new card is the only one, make it default
	var isDefault bool = false
	if len(cards) == 1 && cards[0].Id == paymentMethod.Id {
		err := usersvc.SetDefaultCard(user.Customer.Id, paymentMethod.Id)
		if err != nil {
			reqLogger.WithError(err).Errorln("Unable to set card as default via user-svc")
			WriteError(w, reqLogger, http.StatusInternalServerError)
			return
		}
		isDefault = true
	}

	// Return a cherry-picked version of Stripe's Card struct if the request to
	// the user-svc was successful and a card was returned
	Respond(w, reqLogger, http.StatusOK, exported.Card{
		Id:       paymentMethod.Id,
		Brand:    paymentMethod.Card.Brand,
		ExpMonth: paymentMethod.Card.ExpMonth,
		ExpYear:  paymentMethod.Card.ExpYear,
		LastFour: paymentMethod.Card.Last4,
		Default:  isDefault,
	})
}

// DefaultCard is an HTTP handler for setting a card as default payment method
// for the specified customer in Stripe
func DefaultCard(w http.ResponseWriter, r *http.Request) {
	reqLogger := logging.GetRequestLoggerFromContext(r.Context())
	// Infer customer ID from context
	user, ok := r.Context().Value("User").(structs.UserInternal)
	if !ok {
		reqLogger.Errorln("Unable to get internal user")
		WriteError(w, reqLogger, http.StatusInternalServerError)
		return
	}
	// Set the default card via the user-svc
	paymentMethodId := mux.Vars(r)["card_id"]
	err := usersvc.SetDefaultCard(user.Customer.Id, paymentMethodId)
	if err != nil {
		// Write an appropriate error message if an anticipated error type is found
		switch err {
		case errorspec.CardNotFound:
			WriteError(w, reqLogger, http.StatusBadRequest, "No such card")
		default:
			reqLogger.WithError(err).Errorln("Unable to set card as default via user-svc")
			WriteError(w, reqLogger, http.StatusInternalServerError)
		}
		return
	}

	Respond(w, reqLogger, http.StatusNoContent)
}

// findCard iterates through a list of payment methods and returns a pointer to
// the card with the specified ID. Nil is returned if there is no matching card
func findCard(cardId string, cards []usersvcStructs.PaymentMethod) *usersvcStructs.PaymentMethod {
	for _, card := range cards {
		if card.Id == cardId {
			return &card
		}
	}
	return nil
}

// DeleteCard is an HTTP handler that removes the specified card from Stripe.
// If the card is the default card for the customer, it can only be deleted if
// it is the only card listed for the customer
func DeleteCard(w http.ResponseWriter, r *http.Request) {
	reqLogger := logging.GetRequestLoggerFromContext(r.Context())
	// Infer customer ID from context
	user, ok := r.Context().Value("User").(structs.UserInternal)
	if !ok {
		reqLogger.Errorln("Unable to get internal user")
		WriteError(w, reqLogger, http.StatusInternalServerError)
		return
	}
	// Set the default card via the user-svc
	paymentMethodId := mux.Vars(r)["card_id"]
	// Fetch cards via user-svc
	paymentMethods, err := usersvc.GetCards(user.Customer.Id)
	if err != nil {
		reqLogger.WithError(err).Errorln("Unable to get cards via user-svc")
		WriteError(w, reqLogger, http.StatusInternalServerError)
		return
	}

	// Determine if the specified card belongs to the customer at all
	card := findCard(paymentMethodId, paymentMethods)
	if card == nil {
		WriteError(w, reqLogger,
			http.StatusBadRequest,
			"No such card")
		return
	}

	if card.Default {
		// If the card is default and the only one listed, the customer's default
		// subscription should be canceled as well
		if len(paymentMethods) == 1 {
			// FIXME: Cancel susbscription via user-svc
		} else {
			// Otherwise the neither card nor subscription will be deleted
			WriteError(w, reqLogger,
				http.StatusBadRequest,
				"Can't remove default card unless it's the only one left")
			return
		}
	}
	// Delete the card via user-svc
	if err := usersvc.DeleteCard(user.Customer.Id, card.Id); err != nil {
		// Write an appropriate error message if an anticipated error type is found
		switch err {
		case errorspec.CardNotFound:
			WriteError(w, reqLogger, http.StatusBadRequest, "No such card")
		default:
			reqLogger.WithError(err).Errorln("Unable to delete card via user-svc")
			WriteError(w, reqLogger, http.StatusInternalServerError)
		}
		return
	}

	Respond(w, reqLogger, http.StatusNoContent)
}

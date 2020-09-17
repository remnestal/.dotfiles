package usersvc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	stripeSDK "github.com/stripe/stripe-go"
	"gitlab.com/abios/user-svc/structs"
)

// GetProduct fetches the specified product from the user-svc v2 API
func GetProduct(logger *logrus.Entry, id int) (structs.Product, error) {
	url := fmt.Sprintf("%v/v2/products/default/%v", USERSVC_URL, id)
	// Make GET request
	status, body, err := getRequest(url)
	if err != nil {
		logger.WithError(err).
			WithField("url", url).
			Errorln("Unable to make GET request for default products from user-svc")
		return structs.Product{}, err
	}
	// Ensure that a 200 OK was returned
	if status != http.StatusOK {
		logger.WithError(err).
			WithFields(logrus.Fields{
				"url":    url,
				"status": status,
				"body":   string(body),
			}).Errorln("Product fetch returned non 200 OK response")
		return structs.Product{}, err
	}
	// Unmarshal into the product struct defined by the user-svc
	var prod structs.Product
	if err := json.Unmarshal(body, &prod); err != nil {
		logger.WithError(err).
			WithField("body", body).
			Errorln("Unable to unmarshal product")
		return structs.Product{}, err
	}
	return prod, nil
}

// GetStripePlans fetches the all Stripe product plans from the user-svc v2 API
func GetStripePlans(logger *logrus.Entry) ([]stripeSDK.Plan, error) {
	url := fmt.Sprintf("%v/v2/products/stripe/default", USERSVC_URL)
	// Make GET request
	status, body, err := getRequest(url)
	if err != nil {
		logger.WithError(err).
			WithField("url", url).
			Errorln("Unable to make GET request for Stripe products from user-svc")
		return []structs.Product{}, err
	}
	// Ensure that a 200 OK was returned
	if status != http.StatusOK {
		logger.WithError(err).
			WithFields(logrus.Fields{
				"url":    url,
				"status": status,
				"body":   string(body),
			}).Errorln("Product fetch returned non 200 OK response")
		return []stripeSDK.Plan{}, err
	}
	// Unmarshal into the product struct defined by the user-svc
	var plans []stripeSDK.Plan
	if err := json.Unmarshal(body, &plans); err != nil {
		logger.WithError(err).
			WithField("body", body).
			Errorln("Unable to unmarshal product")
		return []stripeSDK.Plan{}, err
	}
	return plans, nil
}

package handlers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/abios/api-dash-backend/logging"
	"gitlab.com/abios/api-dash-backend/structs/exported"
	"gitlab.com/abios/api-dash-backend/usersvc"
	"gitlab.com/abios/user-svc/structs"
)

// defaultProductIds specifies which products are returned via this HTTP handler
var defaultProductIds []int = []int{2, 3, 4}

// exportDefaultProduct takes a DefaultProduct struct as defined by the user-svc
// and creates an export-friendly representation of it for the dashboard
func exportDefaultProduct(prod structs.Product) exported.DefaultProduct {
	return exported.Product{
		Id:     prod.Id,
		Name:   prod.Name,
		Config: prod.Config,
		Plan: exported.PaymentPlan{
			Amount:    prod.DefaultPlan.Amount,
			Currency:  prod.DefaultPlan.Currency,
			Frequency: prod.DefaultPlan.IntervalCount,
			Interval:  prod.DefaultPlan.Interval,
		},
	}
}

// DefaultProducts is an HTTP handler that returns the list of default products;
// i.e. the products (with default pricing plan) that customers can choose by
// themselves via the dashboard backend
func DefaultProducts(w http.ResponseWriter, r *http.Request) {
	reqLogger := logging.GetRequestLoggerFromContext(r.Context())
	// Fetch each of the listed default products
	products := []exported.Product{}
	for _, id := range defaultProductIds {
		prod, err := usersvc.GetDefaultProduct(reqLogger, id)
		if err != nil {
			reqLogger.WithError(err).
				WithField("product-id", id).
				Errorln("Unable to fetch default product from user-svc")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		products = append(products, exported.DefaultProduct{
			Id: prod.Id,
			Name: prod.Name,
			Config: prod.Config,
		})
	}
	// Prepare response payload
	payload, err := json.Marshal(products)
	if err != nil {
		reqLogger.WithError(err).Errorln("Unable to marshal response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Write response
	if _, err = w.Write(payload); err != nil {
		reqLogger.WithError(err).Errorln("Unable to write response")
		return
	}
}

func StripeProducts(w http.ResponseWriter, r *http.Request) {
  reqLogger := logging.GetRequestLoggerFromContext(r.Context())
  // Fetch all Stripe plans via user-svc
  plans, err := usersvc.GetStripePlans(reqLogger)
  if err != nil {
    reqLogger.WithError(err).Errorln("Unable to fetch Stripe products from user-svc")
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  // Create exportable versions of each product
  expProducts := []exported.Product{}
  for _, plan := range(plans) {
    expProducts = append(expProducts, exportProduct(prod))
  }
  // Prepare response payload
  payload, err := json.Marshal(expProducts)
  if err != nil {
    reqLogger.WithError(err).Errorln("Unable to marshal response")
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  // Write response
  if _, err = w.Write(payload); err != nil {
    reqLogger.WithError(err).Errorln("Unable to write response")
    return
  }
}

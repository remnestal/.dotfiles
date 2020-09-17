package v2api

import (
	"net/http"

	"gitlab.com/abios/user-svc/database/queries"
	"gitlab.com/abios/user-svc/errorspec"
	"gitlab.com/abios/user-svc/httphelper"
	"gitlab.com/abios/user-svc/logging"
	"gitlab.com/abios/user-svc/patch"
	"gitlab.com/abios/user-svc/structs"
	"gitlab.com/abios/user-svc/wraperr"
)

var customerPatchPaths = patch.Paths{
	"/name":            {"replace"},
	"/stripe_id":       {"remove", "replace"},
	"/active_until":    {"remove", "replace"},
	"/payment_source":  {"replace"},
	"/account_manager": {"remove", "replace"},
}

func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	reqLogger := logging.GetRequestLoggerFromContext(r.Context())
	customer, ok := r.Context().Value("customer").(structs.Customer)
	if !ok {
		httphelper.Fail(w, reqLogger,
			http.StatusInternalServerError,
			errorspec.MiddlewareFailure)
		return
	}
	if customer.UpdatedAt != nil {
		if err := httphelper.CompareETag(r, *customer.UpdatedAt); err != nil {
			httphelper.Fail(w, reqLogger,
				wraperr.HttpStatus(err, http.StatusInternalServerError),
				wraperr.Cause(err))
			return
		}
	}
	// Parse the patch parameters
	var args structs.UpdateCustomerParams
	if err := httphelper.ParseBody(r, &args); err != nil {
		httphelper.Fail(w, reqLogger,
			wraperr.HttpStatus(err, http.StatusInternalServerError),
			wraperr.Cause(err))
		return
	}
	// Marshal the Customer-struct and make a comparison
	// customerJSON
	httphelper.Respond(w, reqLogger, http.StatusOK, args)
}

// CustomerInfo is an HTTP handler that accepts a customer ID as URL parameter
// {customer_id} and returns a JSON object with information about that customer
func CustomerInfo(w http.ResponseWriter, r *http.Request) {
	reqLogger := logging.GetRequestLoggerFromContext(r.Context())
	// Extract customer ID
	customer, ok := r.Context().Value("customer").(structs.Customer)
	if !ok {
		httphelper.Fail(w, reqLogger,
			http.StatusInternalServerError,
			errorspec.MiddlewareFailure)
		return
	}
	httphelper.Respond(w, reqLogger, http.StatusOK, customer)
}

// CustomerList is an HTTP handler for listing all customers
func CustomerList(w http.ResponseWriter, r *http.Request) {
	reqLogger := logging.GetRequestLoggerFromContext(r.Context())
	if customers, err := queries.GetCustomers(); err != nil {
		reqLogger.WithError(err).Errorln("Unable to fetch list of customers from database")
		httphelper.Fail(w, reqLogger,
			http.StatusInternalServerError,
			errorspec.DatabaseFailure)
		return
	} else {
		httphelper.Respond(w, reqLogger, http.StatusOK, customers)
	}
}

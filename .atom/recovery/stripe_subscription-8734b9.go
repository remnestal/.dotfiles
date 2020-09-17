package v2api

import (
	"net/http"

	"gitlab.com/abios/user-svc/errorspec"
	"gitlab.com/abios/user-svc/httphelper"
	"gitlab.com/abios/user-svc/logging"
	"gitlab.com/abios/user-svc/database/queries"
)

func GetSubscription(w http.ResponseWriter, r *http.Request) {
	reqLogger := logging.GetRequestLoggerFromContext(r.Context())

	// Extract customer's ID
	customerId, ok := r.Context().Value("customer_id").(int64)
	if !ok {
		httphelper.Fail(w, reqLogger,
			http.StatusInternalServerError,
			errorspec.MiddlewareFailure)
		return
	}

	items, err := queries.GetCustomerProducts(customerId)
	if err != nil {
		httphelper.Fail(w, reqLogger,
			http.StatusInternalServerError,
			errorspec.DatabaseFailure)
		return
	}

	httphelper.Respond(w, reqLogger, http.StatusOK, "test")
}

package v2api

import (
	"net/http"

	stripeSDK "github.com/stripe/stripe-go"

	"gitlab.com/abios/user-svc/httphelper"
	"gitlab.com/abios/user-svc/logging"
	"gitlab.com/abios/user-svc/report"
	"gitlab.com/abios/user-svc/structs"
	"gitlab.com/abios/user-svc/stripe"
)

// ListInvoices is an HTTP handler that fetches all invoices (paid, unpaid, void
// past-due-date, etc) for the specified customer
func ListInvoices(w http.ResponseWriter, r *http.Request) {
	reqLogger := logging.GetRequestLoggerFromContext(r.Context())
	// Extract customer's Stripe ID
	stripeId, ok := r.Context().Value("customer_stripe_id").(string)
	if !ok {
		report.Write(w, reqLogger,
			http.StatusInternalServerError,
			report.ErrorSpec{MiddlewareFailure: true})
		return
	}
	// Fetch all invoices
	params := stripeSDK.InvoiceListParams{}
	params.Filters.AddFilter("customer", "", stripeId)
	invoices := stripe.GetInvoices(&params)

	expInvoices := make(structs.Invoice, len(invoices))
	for _, invoice := range invoices {
		expInvoices = append(expInvoices, structs.Invoice{
			
		})
	}

	httphelper.Respond(w, reqLogger, http.StatusOK, invoices)
}

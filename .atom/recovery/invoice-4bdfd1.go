package structs

type Invoice struct {

}




// InvoiceAmount is a struct for coupling information about the payment in an
// invoice
type InvoiceAmount struct {
	Total     int64  `json:"total"`
	Paid      int64  `json:"paid"`
	Remaining int64  `json:"remaining"`
	Currency  string `json:"currency"`
}

// Invoice is the dashboard representation of a Stripe invoice
type Invoice struct {
	Status           string        `json:"status"`
	Number           string        `json:"number"`
	InvoicePDF       string        `json:"invoice_pdf"`
	InvoiceHosted    string        `json:"invoice_hosted"`
	CollectionMethod string        `json:"collection_method"`
	Created          int64         `json:"created"`
	PeriodStart      int64         `json:"period_start"`
	PeriodEnd        int64         `json:"period_end"`
	DueDate          *int64        `json:"due_date"`
	PaidAt           *int64        `json:"paid_at"`
	Amount           InvoiceAmount `json:"amount"`
}

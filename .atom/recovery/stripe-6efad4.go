package structs

import (
	"github.com/stripe/stripe-go"
)

type StripeTax struct {
	Type   *string                         `json:"type"` // au_abn, eu_vat, in_gst, no_vat, nz_gst
	Value  *string                         `json:"value"`
	Status *stripe.TaxIDVerificationStatus `json:"status"`
}

type StripeBilling struct {
	Email      string    `json:"email"` // Required
	Line1      *string   `json:"line1"` // Required
	Line2      *string   `json:"line2"`
	City       *string   `json:"city"`
	Country    *string   `json:"country"`
	PostalCode *string   `json:"postal_code"`
	State      *string   `json:"state"`
	Tax        StripeTax `json:"tax"`
}

// CardParams specifies what input parameters the user-svc expects in order to
// create a debit/credit card in Stripe
type CardParams struct {
	Number   int64 `json:"number"`
	CVC      int64 `json:"cvc"`
	ExpYear  int64 `json:"exp_year"`
	ExpMonth int64 `json:"exp_month"`
}

type Card struct {
	Fingerprint string `json:"fingerprint"`
	ExpMonth    uint64 `json:"exp_month"`
	ExpYear     uint64 `json:"exp_year"`
	Brand       string `json:"brand"`
	Last4       string `json:"last4"`
}

type PaymentMethod struct {
	Id      string `json:"id"`
	Default bool   `json:"is_default"`
	CardholderName *string `json:"cardholder_name"`
	Card    Card   `json:"card"`
}

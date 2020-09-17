package structs

import (
	jsonpatch "github.com/evanphx/json-patch"
	"gitlab.com/abios/user-svc/errorspec"
)

// Customer is a go struct representation of an Abios customer.
type Customer struct {
	Id               int64   `json:"id"`
	Name             string  `json:"name"`
	StripeId         *string `json:"stripe_id"`
	ActiveUntil      *int64  `json:"active_until"`
	PaymentSource    string  `json:"payment_source"`
	AccountOrigin    string  `json:"account_origin"`
	AccountManagerId *int64  `json:"account_manager_id"`
	UpdatedAt        *int64  `json:"updated_at"`
}

// Parameter-structs for customer & v3-client creation

// checkable defines an interface for shallow attribute validation
type checkable interface {
	Validate() error
}

// evalDependencies invokes the Validate method of each specified dependency
// and forwards any errors encountered. If all dependencies checks out, nil is
// returned instead
func evalDependencies(deps ...checkable) error {
	for _, d := range deps {
		if err := d.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// UserParams defines the parameters needed to specify a new user
type UserParams struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     string  `json:"email"`
}

func (p UserParams) Validate() error {
	/*
		if p.FirstName == "" {
			return errorspec.MissingUserFirstName
		}
		if p.LastName == "" {
			return errorspec.MissingUserLastName
		}
	*/
	if p.Email == "" {
		return errorspec.MissingUserEmail
	}
	return nil
}

// StripeBillingAddressParams defines the parameters needed to set up a new
// Stripe customer's billing address resource
type StripeBillingAddressParams struct {
	Line1      string  `json:"line1"`
	Line2      *string `json:"line2"`
	City       *string `json:"city"`
	PostalCode *string `json:"postal_code"`
	State      *string `json:"state"`
	Country    string  `json:"country"`
}

func (p StripeBillingAddressParams) Validate() error {
	if p.Line1 == "" {
		return errorspec.MissingAddressLine
	}
	if p.Country == "" {
		return errorspec.MissingCountryCode
	}
	return nil
}

// StripeCustomerCreationParams defines the parameters needed to set up a new
// Customer resource in Stripe
type StripeCustomerCreationParams struct {
	BillingEmail   string                     `json:"billing_email"`
	BillingAddress StripeBillingAddressParams `json:"billing_address"`
	Currency       string                     `json:"currency"`
	CouponId       *string                    `json:"coupon_id"`
}

func (p StripeCustomerCreationParams) Validate() error {
	if p.Currency == "" {
		return errorspec.MissingCurrency
	}
	if p.BillingEmail == "" {
		return errorspec.MissingBillingEmail
	}
	return evalDependencies(p.BillingAddress)
}

// StripeCustomerParams defines the representation of a Stripe Customer
// resource; either defined by an Stripe Customer ID or by the parameters
// required to create a new such resource
type StripeCustomerParams struct {
	Id     *string                       `json:"id"`
	Params *StripeCustomerCreationParams `json:"params"`
}

func (p StripeCustomerParams) Validate() error {
	// If Stripe customer ID is set, it can't be empty
	if p.Id != nil && *p.Id == "" {
		return errorspec.MissingStripeId
	}
	// If there's no Stripe Customer ID specified, there must be a Stripe Customer
	// parameter struct present
	if p.Id == nil && p.Params == nil {
		return errorspec.MissingStripeParams
	}
	// Validate Stripe Customer creation parameters if they are set
	if p.Params != nil {
		return evalDependencies(p.Params)
	}
	return nil
}

// CustomerParams defines the parameters required for creating a new customer
// resource in the database and possibly in Stripe
type CustomerParams struct {
	Name           string                `json:"name"`
	User           UserParams            `json:"user"`
	Stripe         *StripeCustomerParams `json:"stripe"`
	ActiveUntil    *int64                `json:"active_until"`
	AccountManager *int64                `json:"account_manager"`
	BillingMethod  string                `json:"billing_method"`
}

func (p CustomerParams) Validate() error {
	if p.Name == "" {
		return errorspec.MissingCustomerName
	}
	// Verify choice of billing method
	switch p.BillingMethod {
	case "stripe-manual", "stripe-auto":
		if p.Stripe == nil {
			return errorspec.MissingStripeResourceDeclaration
		}
	case "other":
	case "":
		return errorspec.MissingBillingMethod
	default:
		return errorspec.InvalidBillingMethod
	}
	// Validate user-params and Stripe-params, if there are any
	if p.Stripe != nil {
		return evalDependencies(p.User, p.Stripe)
	} else {
		return evalDependencies(p.User)
	}
}

// V3ClientConfigParams defines the parameters needed to configure a new v3 API
// client resource
type V3ClientConfigParams struct {
	ReqRateSecond int64 `json:"req_rate_second"`
	HoursBack     int64 `json:"hours_back"`
}

// V3ClientGame defines the parameters required for creating a v3 Client ->
// Game mapping
type V3ClientGame struct {
	GameId          int64  `json:"game_id"`
	PackageId       int64  `json:"package_id"`
	StripeProductId string `json:"stripe_product_id"`
}

func (p V3ClientGame) Validate() error {
	if p.StripeProductId == "" {
		return errorspec.MissingProductId
	}
	return nil
}

// V3ClientParams defines the parameters required to set up a new v3 API client
type V3ClientParams struct {
	CustomerId  int64                `json:"customer_id"`
	ActiveUntil *int64               `json:"active_until"`
	Trialing    bool                 `json:"trialing"`
	Config      V3ClientConfigParams `json:"config"`
	Games       []V3ClientGame       `json:"games"`
}

func (p V3ClientParams) Validate() error {
	for _, game := range p.Games {
		if err := game.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// CreateCustomerParams defines the parameters required to create a new customer
// with associated v2 and v3 clients, as one composite operation
type CreateCustomerParams struct {
	Customer CustomerParams `json:"customer"`
	V3Client V3ClientParams `json:"v3_client"`
}

func (p CreateCustomerParams) Validate() error {
	return evalDependencies(p.Customer, p.V3Client)
}

// UpdateCustomerParams defines the parameters required to update a customer
type UpdateCustomerParams struct {
	Patch jsonpatch.Patch `json:"patch"`
}

// +----------------------------------------------------------------+
// |                                                                |
// |                                                                |
// |    This is a direct port of the user-svc/errorspec package,    |
// |    for compatability-purposes since these changes are not      |
// |    yet in the develop branch of user-svc.                      |
// |                                                                |
// |                                                                |
// +----------------------------------------------------------------+
package errorspec

import (
	"encoding/json"
	"errors"
)

// Error is simply a struct-wrapper for error codes (static string messages).
// A struct is necessary in order to implement the `error` interface, rather
// than just passing strings around
type Error struct {
	Code string `json:"error_code"`
}

// Error returns the error code/message associated with Error struct e
func (e Error) Error() string {
	return e.Code
}

func (e *Error) UnmarshalJSON(data []byte) error {
	var input Error
	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}
	if input.Code == "" {
		return errors.New("Unable to unmarshal Error-spec: error code is not set")
	}
	return nil
}

// The error codes that may be returned by the user-svc
var (
	Unspecified           Error = Error{"unspecified_error"}
	DatabaseFailure       Error = Error{"database_failure"}
	ForeignKeyError       Error = Error{"foreign_key_error"}
	SqlTransactionError   Error = Error{"sql_transaction_error"}
	MiddlewareFailure     Error = Error{"middleware_failure"}
	MalformedRequestBody  Error = Error{"malformed_request_body"}
	MalformedCustomerId   Error = Error{"malformed_customer_id"}
	MalformedProductId    Error = Error{"malformed_product_id"}
	MalformedUserId       Error = Error{"malformed_user_id"}
	InvalidCustomerName   Error = Error{"invalid_customer_name"}
	CustomerNotFound      Error = Error{"customer_not_found"}
	CardNotFound          Error = Error{"card_not_found"}
	DuplicateCard         Error = Error{"duplicate_card"}
	CardExpired           Error = Error{"card_has_expired"}
	InvalidCardCVC        Error = Error{"invalid_card_cvc"}
	InvalidCardNumber     Error = Error{"invalid_card_number"}
	InvalidCardExpYear    Error = Error{"invalid_card_expiry_year"}
	InvalidCardExpMonth   Error = Error{"invalid_card_expiry_month"}
	MissingStripeId       Error = Error{"missing_stripe_id"}
	MissingStripeResource Error = Error{"missing_stripe_resource"}
)

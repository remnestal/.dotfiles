package stripe

import (
	"github.com/stripe/stripe-go"
	"gitlab.com/abios/user-svc/errorspec"
)

type ErrorMapping struct {
	Critical bool
	Code     errorspec.Error
}

// errorMap represents the mapping between Stripe's error codes and the error
// codes that are returned by the user-svc. This mapping is surjective
var errorMap map[stripe.ErrorCode]errorspec.Error = map[stripe.ErrorCode]ErrorMapping{
	stripe.ErrorCodeResourceMissing: ErrorMapping{
		Code:     errorspec.MissingStripeResource,
		Critical: true,
	},
	stripe.ErrorCodeExpiredCard: ErrorMapping{
		Code:     errorspec.CardExpired,
		Critical: false,
	},
	stripe.ErrorCodeIncorrectCVC: ErrorMapping{
		Code:     errorspec.InvalidCardCVC,
		Critical: false,
	},
	stripe.ErrorCodeIncorrectNumber: ErrorMapping{
		Code:     errorspec.InvalidCardNumber,
		Critical: false,
	},
	stripe.ErrorCodeInvalidExpiryMonth: ErrorMapping{
		Code:     errorspec.InvalidCardExpMonth,
		Critical: false,
	},
	stripe.ErrorCodeInvalidExpiryYear: ErrorMapping{
		Code:     errorspec.InvalidCardExpYear,
		Critical: false,
	},
	stripe.ErrorCodeInvalidNumber: ErrorMapping{
		Code:     errorspec.InvalidCardNumber,
		Critical: false,
	},
	stripe.ErrorCodeInvalidCVC: ErrorMapping{
		Code:     errorspec.InvalidCardCVC,
		Critical: false,
	},
}

// ParseError accepts an error instance and tries to parse it as if it was an
// error returned by the Stripe SDK. If this is the case, the appropriate
// internal error code is returned, or an error of type "unmapped stripe error"
// if there is no mapping for that Stripe error. If the passed error was not
// issued by Stripe, it is simply returned as an unspecified error
func ParseError(err error) errorspec.Error {
	if stripeErr, ok := err.(*stripe.Error); ok {
		if code, exists := errorMap[stripeErr.Code]; exists {
			return code
		} else {
			stripeLogger.WithField("code", stripeErr.Code).Errorln("Encountered an unmapped Stripe error code")
			return errorspec.UnMappedStripeError
		}
	} else {
		stripeLogger.WithError(err).Errorln("Tried to parse non-Stripe error")
		return errorspec.Unknown
	}
}

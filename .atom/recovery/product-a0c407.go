package structs

import (
  "github.com/stripe/stripe-go"
)

type ProductConfig struct {
  TimeAccess struct {
    HoursBack int64 `json:"hours_back"`
    HoursForward int64 `json:"hours_forward"`
  } `json:"time_access"`
  RequestRates struct {
    Minute int64 `json:"minute"`
    Search int64 `json:"search"`
    Second int64 `json:"second"`
  } `json:"request_rates"`
}

type Product struct {
  Id int64 `json:"id"`
  Name string `json:"name"`
  Config ProductConfig `json:"config"`
  StripeId string `json:"stripe_id"`
  DefaultPlan stripe.Plan
}

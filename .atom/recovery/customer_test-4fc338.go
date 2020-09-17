package v2api

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	. "github.com/steinfletcher/apitest-jsonpath"

	"gitlab.com/abios/user-svc/server"
)

func TestGetCustomer(t *testing.T) {
	t.Run("Fetch dummy customer 1", func(t *testing.T) {
		apitest.New().
			Handler(server.Routes()).
			Get(fmt.Sprintf("/v2/customer/%v", CUSTOMER_1_ID)).
			Expect(t).
			Status(http.StatusOK).
			Assert(Equal(`$.id`, number(CUSTOMER_1_ID))).
			Assert(Equal(`$.name`, "abios-user-svc-test-1")).
			Assert(Equal(`$.stripe_id`, "cus_G5m8Rw09eJkNnj")).
			Assert(Equal(`$.active_until`, number(1570699577))).
			Assert(Equal(`$.payment_source`, "stripe-auto")).
			Assert(Equal(`$.account_origin`, "managed")).
			Assert(Equal(`$.account_manager_id`, nil)).
			Assert(Present(`$.updated_at`)).
			Assert(Present(`$.updated_at`)).
			End()
	})
	t.Run("Fetch dummy customer 2", func(t *testing.T) {
		apitest.New().
			Handler(server.Routes()).
			Get(fmt.Sprintf("/v2/customer/%v", CUSTOMER_2_ID)).
			Expect(t).
			Status(http.StatusOK).
			Assert(Equal(`$.id`, number(CUSTOMER_2_ID))).
			Assert(Equal(`$.name`, "abios-user-svc-test-2")).
			Assert(Equal(`$.stripe_id`, nil)).
			Assert(Equal(`$.active_until`, number(1570699577))).
			Assert(Equal(`$.payment_source`, "other")).
			Assert(Equal(`$.account_origin`, "managed")).
			Assert(Equal(`$.account_manager_id`, nil)).
			Assert(Present(`$.created_at`)).
			Assert(Equal(`$.updated_at`, nil)).
			End()
	})
	t.Run("Fetch dummy customer 3", func(t *testing.T) {
		apitest.New().
			Handler(server.Routes()).
			Get(fmt.Sprintf("/v2/customer/%v", CUSTOMER_3_ID)).
			Expect(t).
			Status(http.StatusOK).
			Assert(Equal(`$.id`, number(CUSTOMER_3_ID))).
			Assert(Equal(`$.name`, "abios-user-svc-test-3")).
			Assert(Equal(`$.stripe_id`, "cus_Fwr7sMMns613n5")).
			Assert(Equal(`$.active_until`, number(1570699577))).
			Assert(Equal(`$.payment_source`, "stripe-manual")).
			Assert(Equal(`$.account_origin`, "managed")).
			Assert(Equal(`$.account_manager_id`, number(262))).
			Assert(Present(`$.created_at`)).
			Assert(Present(`$.updated_at`)).
			End()
	})
	// FIXME: every customer object should adhere to the OAS spec
}

func TestGetCustomerMalformedId(t *testing.T) {
	MalformedCustomerId(t, "/v2/customer/%v")
}

func TestGetCustomerInvalidId(t *testing.T) {
	InvalidCustomerId(t, "/v2/customer/%v")
}

func InvalidCustomerId(t *testing.T, url string) {
	t.Run("Non-existing negative customer-ID", func(t *testing.T) {
		apitest.New().
			Handler(server.Routes()).
			Get(fmt.Sprintf(url, -1)).
			Expect(t).
			Status(http.StatusNotFound).
			Assert(Equal(`$.error_code`, "customer_not_found")).
			End()
	})
	t.Run("Non-existing positive customer-ID", func(t *testing.T) {
		apitest.New().
			Handler(server.Routes()).
			Get(fmt.Sprintf(url, maxInt())).
			Expect(t).
			Status(http.StatusNotFound).
			Assert(Equal(`$.error_code`, "customer_not_found")).
			End()
	})
	// FIXME: every error message should adhere to the OAS spec
}

func MalformedCustomerId(t *testing.T, url string) {
	t.Run("Non integer ID", func(t *testing.T) {
		apitest.New().
			Handler(server.Routes()).
			Get(fmt.Sprintf(url, "abc")).
			Expect(t).
			Status(http.StatusBadRequest).
			Assert(Equal(`$.error_code`, "malformed_customer_id")).
			End()
	})
	// FIXME: every error message should adhere to the OAS spec
}

package paymentpage

import (
	"reflect"
	"testing"
)

func TestPayment(t *testing.T) {
	t.Parallel()
	payment := NewPayment(11, "test_payment_id")

	compareMap := map[string]interface{}{"project_id": 11, "payment_id": "test_payment_id"}
	equal := reflect.DeepEqual(compareMap, payment.GetParams())

	if !equal {
		t.Error(
			"For", "NewPayment",
			"expected", compareMap,
			"got", payment.GetParams(),
		)
	}
}

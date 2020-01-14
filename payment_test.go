package paymentpage

import (
	"reflect"
	"testing"
	"time"
)

func TestPayment(t *testing.T) {
	t.Parallel()
	timeNow, _ := time.Parse(time.RFC3339, "2222-01-01T11:11:11Z")
	payment := NewPayment(11, nil)
	payment.SetParam(ParamPaymentId, "test_payment_id")
	payment.SetParam(ParamBestBefore, timeNow)

	compareMap := map[string]interface{}{"project_id": 11, "payment_id": "test_payment_id", "best_before": "2222-01-01T11:11:11Z", "interface_type": `{"id": 20}`}
	equal := reflect.DeepEqual(compareMap, payment.GetParams())

	if !equal {
		t.Error(
			"For", "NewPayment",
			"expected", compareMap,
			"got", payment.GetParams(),
		)
	}
}

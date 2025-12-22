package paymentpage

import (
	"encoding/base64"
	"encoding/json"
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
	payment.SetParam(ParamOperationType, "sale")

	compareMap := map[string]interface{}{"project_id": 11, "payment_id": "test_payment_id", "best_before": "2222-01-01T11:11:11Z", "interface_type": `{"id": 20}`, "operation_type": "sale"}
	equal := reflect.DeepEqual(compareMap, payment.GetParams())

	if !equal {
		t.Error(
			"For", "NewPayment",
			"expected", compareMap,
			"got", payment.GetParams(),
		)
	}
}

func TestCardOperationType(t *testing.T) {
	t.Parallel()
	payment := NewPayment(11, nil)
	payment.SetParam(ParamCardOperationType, "sale")

	compareMap := map[string]interface{}{"project_id": 11, "interface_type": `{"id": 20}`, "operation_type": "sale"}
	equal := reflect.DeepEqual(compareMap, payment.GetParams())

	if !equal {
		t.Error(
			"For", "NewPayment",
			"expected", compareMap,
			"got", payment.GetParams(),
		)
	}
}

func TestSetBookingInfo(t *testing.T) {
	expected := LoadJsonFromFile(t, "booking_info.json")
	payment := NewPayment(11, nil).SetBookingInfo(expected)

	actualBookingInfoJson, _ := base64.StdEncoding.DecodeString(payment.GetParams()["booking_info"].(string))

	var actual map[string]any 
	_ = json.Unmarshal(actualBookingInfoJson, &actual)

	equal := reflect.DeepEqual(expected, actual)
	if !equal {
		t.Error(
			"For", "NewPayment",
			"expected", expected,
			"got", actual,
		)
	}
}

func TestSetBookingInfoNilPanics(t *testing.T) {
	payment := NewPayment(11, nil)
	defer func() {
		if recover() == nil {
			t.Error("expected panic for nil bookingInfo")
		}
	}()
	payment.SetBookingInfo(nil)
}
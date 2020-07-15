package paymentpage

import (
	"reflect"
	"testing"
)

func TestCallback(t *testing.T) {
	t.Parallel()

	paymentId := "11"
	float64PaymentId := "11111111111111.1"
	paymentStatus := "success"
	signature := "Rk9xmCfiv/BJbVrCz+oazsOOuiMqrktLVXruRLM9WJ2zmYvufwOS7uxz5Pd36kfKPqbBwKZjMs/EEzF/VsbbNw=="
	callbackDataGeneralSign := `{
		"general": {
			"signature": "` + signature + `"
		},
		"payment": {
			"id": "` + paymentId + `",
			"status": "` + paymentStatus + `"
		},
		"provider_extra_fields": {
			"extended": {}
		}
	}`
	callbackData := `{
		"payment": {
			"id": "` + paymentId + `",
			"status": "` + paymentStatus + `"
		},
		"signature": "` + signature + `"
	}`
	callbackDataPaymentInt := `{
		"payment": {
			"id": "` + paymentId + `",
			"status": "` + paymentStatus + `"
		},
		"signature": "` + signature + `"
	}`
	callbackDataInvalidSign := `{
		"payment": {
			"id": "` + paymentId + `"
		},
		"signature": "f2g3h4j5"
	}`
	callbackDataWithFloatId := `{
		"payment": {
			"id": 11111111111111.1
		}
	}`
	recurringUpdateCallbackData := `{
		"project_id": 123,
		"recurring": {
			"id": 321,
			"status": "active",
			"type": "Y",
			"currency": "EUR",
			"exp_year": "2025",
			"exp_month": "12",
			"period": "D",
			"time": "11:00"
		}
	}`

	signatureHandler := NewSignatureHandler("qwerty")
	callback, err := NewCallback(*signatureHandler, callbackData)

	if err != nil {
		t.Error(
			"For", "NewCallback",
			"expected", "Callback",
			"got", err.Error(),
		)
	}

	comparePayment := map[string]interface{}{
		"id":     paymentId,
		"status": paymentStatus,
	}

	equal := reflect.DeepEqual(comparePayment, callback.GetPayment())

	if !equal {
		t.Error(
			"For", "GetPayment",
			"expected", comparePayment,
			"got", callback.GetPayment(),
		)
	}

	if callback.GetPaymentStatus() != paymentStatus {
		t.Error(
			"For", "GetPaymentStatus",
			"expected", paymentStatus,
			"got", callback.GetPaymentStatus(),
		)
	}

	if callback.GetPaymentId() != paymentId {
		t.Error(
			"For", "GetPaymentId",
			"expected", paymentId,
			"got", callback.GetPaymentId(),
		)
	}

	cbSign, _ := callback.getSignature()

	if cbSign != signature {
		t.Error(
			"For", "getSignature",
			"expected", signature,
			"got", cbSign,
		)
	}

	callbackIntPayment, _ := NewCallback(*signatureHandler, callbackDataPaymentInt)
	emptyValue := callbackIntPayment.GetParam("undefined_key")

	if emptyValue != nil {
		t.Error(
			"For", "getParamByName",
			"expected", nil,
			"got", emptyValue,
		)
	}

	if callbackIntPayment.GetPaymentId() != paymentId {
		t.Error(
			"For", "GetPaymentId",
			"expected", paymentId,
			"got", callbackIntPayment.GetPaymentId(),
		)
	}

	callbackInvalidData, err := NewCallback(*signatureHandler, "{sfsggsg")
	_ = callbackInvalidData

	if err.Error() != "invalid callback data" {
		t.Error(
			"For", "NewCallback",
			"expected", "invalid callback data",
			"got", "Callback",
		)
	}

	callbackInvalidSignature, err := NewCallback(*signatureHandler, callbackDataInvalidSign)
	_ = callbackInvalidSignature

	if err.Error() != "invalid signature" {
		t.Error(
			"For", "NewCallback",
			"expected", "invalid signature",
			"got", "Callback",
		)
	}

	callbackInvalidGeneralSign, err := NewCallback(*signatureHandler, callbackDataGeneralSign)
	_ = callbackInvalidGeneralSign

	if err != nil && err.Error() != "invalid signature" {
		t.Error(
			"For", "NewCallback",
			"expected", "invalid signature",
			"got", "Callback",
		)
	}

	callbackWithFloatId, _ := NewCallback(*signatureHandler, callbackDataWithFloatId)

	if callbackWithFloatId.GetPaymentId() != float64PaymentId {
		t.Error(
			"For", "NewCallback",
			"expected", float64PaymentId,
			"got", callbackWithFloatId.GetPaymentId(),
		)
	}

	callbackWithoutPayment, err := NewCallback(*signatureHandler, recurringUpdateCallbackData)
	payment := callbackWithoutPayment.GetPayment()

	if payment != nil {
		t.Error(
			"For", "NewCallback",
			"expected", nil,
			"got", payment,
		)
	}

	if err == nil || err.Error() != "signature undefined" {
		sign, _ := callbackWithoutPayment.getSignature()
		t.Error(
			"For", "NewCallback",
			"expected", "signature undefined",
			"got", sign,
		)
	}
}

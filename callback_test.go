package paymentpage

import (
	"reflect"
	"testing"
)

func TestCallback(t *testing.T) {
	t.Parallel()

	paymentId := "11"
	paymentStatus := "success"
	signature := "k7XHz9UctgxYDmbjSjl7xi7PPl59hX1yuLQG6MH7sC5I4HiMhqyNU89UXkX7tM9bIH/TSXwlJkeH0qvVx8i+hA=="
	callbackData := "{" +
		"\"body\": {" +
		"\"payment\": {" +
		"\"id\": \"" + paymentId + "\"," +
		"\"status\": \"" + paymentStatus + "\"" +
		"}," +
		"\"signature\": \"" + signature + "\"" +
		"}" +
		"}"
	callbackDataPaymentInt := "{" +
		"\"body\": {" +
		"\"payment\": {" +
		"\"id\": " + paymentId + "," +
		"\"status\": \"" + paymentStatus + "\"" +
		"}," +
		"\"signature\": \"" + signature + "\"" +
		"}" +
		"}"
	callbackDataInvalidSign := "{" +
		"\"body\": {" +
		"\"payment\": {" +
		"\"id\": \"" + paymentId + "\"" +
		"}," +
		"\"signature\": \"f2g3h4j5\"" +
		"}" +
		"}"

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
			"expected", paymentStatus,
			"got", callback.GetPaymentId(),
		)
	}

	if callback.getSignature() != signature {
		t.Error(
			"For", "getSignature",
			"expected", signature,
			"got", callback.getSignature(),
		)
	}

	callbackIntPayment, _ := NewCallback(*signatureHandler, callbackDataPaymentInt)
	someData := map[string]interface{}{"qwerty": "111"}
	emptyValue := callbackIntPayment.getParamByName("undefined_key", someData)

	if emptyValue != "" {
		t.Error(
			"For", "getParamByName",
			"expected", "",
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
}

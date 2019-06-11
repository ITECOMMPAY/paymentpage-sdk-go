package paymentpage

import (
	"reflect"
	"testing"
)

func TestCallback(t *testing.T) {
	t.Parallel()

	paymentId := "11"
	paymentStatus := "success"
	signature := "Rk9xmCfiv/BJbVrCz+oazsOOuiMqrktLVXruRLM9WJ2zmYvufwOS7uxz5Pd36kfKPqbBwKZjMs/EEzF/VsbbNw=="
	callbackDataRecursive := "{" +
		"\"body\": {" +
		"\"payment\": {" +
		"\"id\": \"" + paymentId + "\"," +
		"\"status\": \"" + paymentStatus + "\"" +
		"}," +
		"\"signature\": \"" + signature + "\"" +
		"}" +
		"}"
	callbackData := "{" +
		"\"payment\": {" +
		"\"id\": \"" + paymentId + "\"," +
		"\"status\": \"" + paymentStatus + "\"" +
		"}," +
		"\"signature\": \"" + signature + "\"" +
		"}"
	callbackDataPaymentInt := "{" +
		"\"payment\": {" +
		"\"id\": " + paymentId + "," +
		"\"status\": \"" + paymentStatus + "\"" +
		"}," +
		"\"signature\": \"" + signature + "\"" +
		"}"
	callbackDataInvalidSign := "{" +
		"\"payment\": {" +
		"\"id\": \"" + paymentId + "\"" +
		"}," +
		"\"signature\": \"f2g3h4j5\"" +
		"}"
	callbackDataWithSlice := "{" +
		"\"errors\": [" +
		"{" +
		"\"message\": \"error1\"," +
		"\"code\": 1," +
		"\"fail\": true" +
		"}," +
		"{" +
		"\"message\": \"error2\"," +
		"\"code\": 2," +
		"\"fail\": false" +
		"}," +
		"{" +
		"\"message\": \"error3\"," +
		"\"code\": 3," +
		"\"fail\": {" +
		"\"sub-fail\": {" +
		"\"test\": 1" +
		"}" +
		"}" +
		"}" +
		"]," +
		"\"sum_converted\": {" +
		"\"amount\": 10000," +
		"\"currency\": \"EUR\"" +
		"}," +
		"\"signature\": \"oXF9a0F80FBkT5plV1aVcQIVSWr3l07StkQ2izUKmy//H2S9gMX982Kgm4tXB4+Ze1S5E1jeKhwheIgYMZ4J+w==\"" +
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

	callbackInvalidSignatureRecursive, err := NewCallback(*signatureHandler, callbackDataRecursive)
	_ = callbackInvalidSignatureRecursive
	if err.Error() != "invalid signature" {
		t.Error(
			"For", "NewCallback",
			"expected", "invalid signature",
			"got", "Callback",
		)
	}

	callbackWithSlice, err := NewCallback(*signatureHandler, callbackDataWithSlice)
	_ = callbackWithSlice

	if err != nil {
		t.Error(
			"For", "NewCallback",
			"expected", "Valid callback with slice",
			"got", err.Error(),
		)
	}
}

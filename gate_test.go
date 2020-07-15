package paymentpage

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGate(t *testing.T) {
	t.Parallel()

	paymentId := "11"
	paymentStatus := "success"
	signature := "Rk9xmCfiv/BJbVrCz+oazsOOuiMqrktLVXruRLM9WJ2zmYvufwOS7uxz5Pd36kfKPqbBwKZjMs/EEzF/VsbbNw=="
	callbackData := `{
		"payment": {
			"id": "` + paymentId + `",
			"status": "` + paymentStatus + `"
		},
		"signature": "` + signature + `"
	}`

	comparePaymentHost := "test.test"
	comparePaymentPath := "/pay"
	comparePaymentQuery := map[string]interface{}{
		"project_id":             "11",
		"payment_id":             "test_payment_id",
		"interface_type":         `{"id": 20}`,
		"payment_currency":       "EUR",
		"payment_amount":         "1000",
		"some_future_bool_param": "true",
		"signature":              "+NGChjO/3L6vSJkEUSXBDPJBSUuEu4rXw4wtAoXiTDATSMerNixVYKdh9Cg2jTXSu1Ez9R+LxX/ioWr70Tlxew==",
	}

	payment := NewPayment(11, "test_payment_id")
	payment.SetParam("payment_currency", "EUR")
	payment.SetParam("payment_amount", 1000)
	payment.SetParam("some_future_bool_param", true)

	gate := NewGate("qwerty")

	gate.SetBaseUrl("http://test.test/pay")
	paymentPageUrl := gate.GetPaymentPageUrl(*payment)
	parsedUrl, err := url.Parse(paymentPageUrl)
	query, _ := url.ParseQuery(parsedUrl.RawQuery)

	if err != nil {
		t.Error(
			"For", "GetPaymentPageUrl",
			"expected", "valid URL",
			"got", paymentPageUrl,
		)
	}

	if parsedUrl.Host != comparePaymentHost {
		t.Error(
			"For", "GetPaymentPageUrl", payment,
			"expected", comparePaymentHost,
			"got", parsedUrl.Host,
		)
	}

	if parsedUrl.Path != comparePaymentPath {
		t.Error(
			"For", "GetPaymentPageUrl", payment,
			"expected", comparePaymentPath,
			"got", parsedUrl.Path,
		)
	}

	realQuery := map[string]interface{}{}

	for key, value := range query {
		realQuery[key] = value[0]
	}

	equal := reflect.DeepEqual(comparePaymentQuery, realQuery)

	if !equal {
		t.Error(
			"For", "GetPaymentPageUrl",
			"expected", comparePaymentQuery,
			"got", realQuery,
		)
	}

	callback, err := gate.HandleCallback(callbackData)
	signatureHandler := NewSignatureHandler("qwerty")
	callbackTest, errTest := NewCallback(*signatureHandler, callbackData)

	if reflect.TypeOf(callback) != reflect.TypeOf(callbackTest) || err != errTest {
		t.Error(
			"For", "HandleCallback",
			"expected", callbackTest, errTest,
			"got", callback, err,
		)
	}
}

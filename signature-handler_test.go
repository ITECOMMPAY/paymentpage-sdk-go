package paymentpage

import (
	"encoding/json"
	"testing"
)

func TestSignatureHandler(t *testing.T) {
	t.Parallel()

	twoLevelsParam := map[string]interface{}{
		"param1": 1, "param2": 2,
	}
	params := map[string]interface{}{
		"project_id":       11,
		"payment_currency": "EUR",
		"frame_mode":       "popup",
		"some_bool_param":  true,
		"two_levels_param": twoLevelsParam,
	}
	signatureHandler := NewSignatureHandler("qwerty")

	if !signatureHandler.Check("ehwG+e+Sq6b0QEek2hmHGUFzViX9/Xfcgnh0Ysn1+AqCY8BsKhC/Vcyr4XH8CKXW3iS0vW9S2mTwOt/0Ozf5pg==", params) {
		t.Error(
			"For", "qwerty", params,
			"expected", true,
			"got", false,
		)
	}

	if signatureHandler.Check("tIf3gh4jZHReVS5aOz24/VfAorw7mQ==", params) {
		t.Error(
			"For", "qwerty", params,
			"expected", false,
			"got", true,
		)
	}

	signatureHandler.SetSort(false)

	if signatureHandler.sort {
		t.Error(
			"For", "SetSort", false,
			"expected", false,
			"got", true,
		)
	}

	jsonParams := `{"a": {"b": 1111111111111111, "c": 2 }, "d": false, "f": [ "g", { "h": 7, "i": { "k": 8 } } ], "e": [ 4, 5 ], "j": null}`
	expectedString := "a:b:1111111111111111;a:c:2;d:0;e:0:4;e:1:5;f:0:g;f:1:h:7;f:1:i:k:8;j:"
	parsedParams := make(map[string]interface{})
	parseError := json.Unmarshal([]byte(jsonParams), &parsedParams)
	_ = parseError
	signatureHandler.SetSort(true)
	stringParamsToSign := signatureHandler.getStringParamsToSign(parsedParams)

	if stringParamsToSign != expectedString {
		t.Error(
			"For", jsonParams,
			"expected", expectedString,
			"got", stringParamsToSign,
		)
	}
}

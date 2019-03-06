package paymentpage

import (
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
}

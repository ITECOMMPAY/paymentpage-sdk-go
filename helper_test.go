package paymentpage

import (
	"testing"
)

func TestHelper(t *testing.T) {
	if concat("", "") != "" {
		t.Error(
			"For", "concat",
			"expected", "",
			"got", concat("", ""),
		)
	}

	if getStringBool(false) != "0" {
		t.Error(
			"For", "getStringBool",
			"expected", "0",
			"got", getStringBool(false),
		)
	}
}

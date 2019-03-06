package paymentpage

import (
	"net/url"
	"strconv"
	"strings"
)

// Structure for build payment URL
type PaymentPage struct {
	// Payment domain with path
	baseUrl string

	// Signature handler for generate signature
	signatureHandler SignatureHandler
}

// Method for set base payment page URL
func (p *PaymentPage) SetBaseUrl(baseUrl string) *PaymentPage {
	p.baseUrl = baseUrl

	return p
}

// Method build payment URL
func (p *PaymentPage) GetUrl(payment Payment) string {
	signature := p.signatureHandler.Sign(payment.GetParams())

	queryArray := []string{}

	for key, value := range payment.GetParams() {
		preparedValue := ""

		switch value := value.(type) {
		case string:
			preparedValue = value
		case int:
			preparedValue = strconv.Itoa(value)
		case bool:
			preparedValue = strconv.FormatBool(value)
		}

		queryArray = append(queryArray, concat(concat(key, "="), url.QueryEscape(preparedValue)))
	}

	queryString := strings.Join(queryArray, "&")
	queryString = concat(queryString, concat("&signature=", url.QueryEscape(signature)))

	return concat(p.baseUrl, concat("?", queryString))
}

// Constructor for PaymentPage structure
func NewPaymentPage(signatureHandler SignatureHandler) *PaymentPage {
	paymentPage := PaymentPage{"https://paymentpage.ecommpay.com/payment", signatureHandler}

	return &paymentPage
}

package paymentpage

import (
	"net/url"
	"strconv"
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
	link, err := url.Parse(p.baseUrl)
	if err != nil {
		panic(err)
	}

	query := link.Query()

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

		query.Set(key, preparedValue)
	}
	
	signature := p.signatureHandler.Sign(payment.GetParams())
	query.Set("signature", signature)

	link.RawQuery = query.Encode()

	return link.String()
}

// Constructor for PaymentPage structure
func NewPaymentPage(signatureHandler SignatureHandler) *PaymentPage {
	paymentPage := PaymentPage{"https://paymentpage.ecommpay.com/payment", signatureHandler}

	return &paymentPage
}

package paymentpage

import (
	"net/url"
	"strconv"
	"strings"
)

type PaymentPage struct {
	baseUrl          string
	signatureHandler SignatureHandler
}

func (p *PaymentPage) SetBaseUrl(baseUrl string) *PaymentPage {
	p.baseUrl = baseUrl

	return p
}

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

func NewPaymentPage(signatureHandler SignatureHandler) *PaymentPage {
	paymentPage := PaymentPage{"https://paymentpage.ecommpay.com/payment", signatureHandler}

	return &paymentPage
}

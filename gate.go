package paymentpage

type Gate struct {
	signatureHandler SignatureHandler
	paymentPage      PaymentPage
}

func (g *Gate) SetBaseUrl(url string) *Gate {
	g.paymentPage.SetBaseUrl(url)

	return g
}

func (g *Gate) GetPaymentPageUrl(payment Payment) string {
	return g.paymentPage.GetUrl(payment)
}

func (g *Gate) HandleCallback(callbackData string) (*Callback, error) {
	callback, callbackError := NewCallback(g.signatureHandler, callbackData)

	return callback, callbackError
}

func NewGate(secret string) *Gate {
	signatureHandler := NewSignatureHandler(secret)
	paymentPage := NewPaymentPage(*signatureHandler)
	gate := Gate{*signatureHandler, *paymentPage}

	return &gate
}

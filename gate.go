package paymentpage

// Structure for communicate with our
type Gate struct {
	// Instance for check signature
	signatureHandler SignatureHandler

	// Instance for build payment URL
	paymentPage PaymentPage
}

// Method for set base payment page URL
func (g *Gate) SetBaseUrl(url string) *Gate {
	g.paymentPage.SetBaseUrl(url)

	return g
}

// Method build payment URL
func (g *Gate) GetPaymentPageUrl(payment Payment) string {
	return g.paymentPage.GetUrl(payment)
}

// Method for handling callback
func (g *Gate) HandleCallback(callbackData string) (*Callback, error) {
	callback, callbackError := NewCallback(g.signatureHandler, callbackData)

	return callback, callbackError
}

// Constructor for Gate structure
func NewGate(secret string) *Gate {
	signatureHandler := NewSignatureHandler(secret)
	paymentPage := NewPaymentPage(*signatureHandler)
	gate := Gate{*signatureHandler, *paymentPage}

	return &gate
}

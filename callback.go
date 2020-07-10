package paymentpage

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

// Constants with possible statuses of payment
const (
	PaymentStatusSuccess           string = "success"
	PaymentStatusDecline           string = "decline"
	PaymentStatusAW3DS             string = "awaiting 3ds result"
	PaymentStatusAWRedirect        string = "awaiting redirect result"
	PaymentStatusAWCustomer        string = "awaiting customer"
	PaymentStatusAWClarification   string = "awaiting clarification"
	PaymentStatusAWCapture         string = "awaiting capture"
	PaymentStatusCancelled         string = "cancelled"
	PaymentStatusRefunded          string = "refunded"
	PaymentStatusPartiallyRefunded string = "partially refunded"
	PaymentStatusProcessing        string = "processing"
	PaymentStatusError             string = "error"
	PaymentStatusReversed          string = "reversed"
)

// Structure for processing callbacks
type Callback struct {
	// Instance for check callback signature
	signatureHandler SignatureHandler

	// Raw callback data
	callbackData string

	// Decoded callback data
	parsedData map[string]interface{}

	// Callback signature
	signature string
}

// Return map with payment data
func (c *Callback) GetPayment() interface{} {
	return c.GetParam("payment")
}

// Return payment status
func (c *Callback) GetPaymentStatus() interface{} {
	status := c.GetParam("payment.status")

	return status.(string)
}

// Return our payment id
func (c *Callback) GetPaymentId() interface{} {
	id := c.GetParam("payment.id")

	switch id := id.(type) {
	case float64:
		return strconv.FormatFloat(id, 'f', -1, 64)
	default:
		return id.(string)
	}
}

// Get callback param by path name
func (c *Callback) GetParam(pathStr string) interface{} {
	path := strings.Split(pathStr, ".")
	cbData := c.parsedData
	var value interface{}

	for _, key := range path {
		tmpVal, find := cbData[key]
		value = tmpVal

		if !find {
			break
		}

		switch value := value.(type) {
		case map[string]interface{}:
			cbData = value
		}
	}

	return value
}

// Return callback signature
func (c *Callback) getSignature() (string, error) {
	if c.signature != "" {
		return c.signature, nil
	}

	signPaths := []string{
		"signature",
		"general.signature",
	}

	for _, signPath := range signPaths {
		sign := c.GetParam(signPath)

		if sign != nil {
			c.signature = sign.(string)

			return c.signature, nil
		}
	}

	return "", errors.New("signature undefined")
}

// Check that signature is valid
func (c *Callback) checkSignature() error {
	signature, err := c.getSignature()

	if err != nil {
		return err
	}

	c.removeParam("signature", c.parsedData)

	if !c.signatureHandler.Check(signature, c.parsedData) {
		return errors.New("invalid signature")
	}

	return nil
}

// Method for remove value in multilevel map
func (c *Callback) removeParam(name string, data map[string]interface{}) {
	if _, find := data[name]; find {
		delete(data, name)
	}

	for _, value := range data {
		switch value := value.(type) {
		case map[string]interface{}:
			c.removeParam(name, value)
		}
	}
}

// Method for decode callback data
func (c *Callback) parseCallbackData() error {
	parseError := json.Unmarshal([]byte(c.callbackData), &c.parsedData)

	if parseError != nil {
		return errors.New("invalid callback data")
	}

	return nil
}

// Constructor for Callback structure
func NewCallback(signatureHandler SignatureHandler, callbackData string) (*Callback, error) {
	callback := Callback{}
	callback.signatureHandler = signatureHandler
	callback.callbackData = callbackData

	if parseError := callback.parseCallbackData(); parseError != nil {
		return &callback, parseError
	}

	if signatureError := callback.checkSignature(); signatureError != nil {
		return &callback, signatureError
	}

	return &callback, nil
}

package paymentpage

import (
	"encoding/json"
	"errors"
	"strconv"
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
	return c.getParamByName("payment", c.parsedData)
}

// Return payment status
func (c *Callback) GetPaymentStatus() string {
	status := c.getParamByName("status", c.parsedData)

	return status.(string)
}

// Return our payment id
func (c *Callback) GetPaymentId() string {
	id := c.getParamByName("id", c.parsedData)

	switch id := id.(type) {
	case float64:
		return strconv.FormatFloat(id, 'f', 0, 64)
	default:
		return id.(string)
	}
}

// Return callback signature
func (c *Callback) getSignature() string {
	if c.signature == "" {
		c.signature = c.getParamByName("signature", c.parsedData).(string)
	}

	return c.signature
}

// Check that signature is valid
func (c *Callback) checkSignature() error {
	signature := c.getSignature()
	c.removeParam("signature", c.parsedData)

	if !c.signatureHandler.Check(signature, c.parsedData) {
		return errors.New("invalid signature")
	}

	return nil
}

// Method for get value in multilevel map by key
func (c *Callback) getParamByName(name string, data map[string]interface{}) interface{} {
	if value, find := data[name]; find {
		return value
	}

	for _, value := range data {
		switch value := value.(type) {
		case map[string]interface{}:
			param := c.getParamByName(name, value)

			if param != "" {
				return param
			}
		}
	}

	return ""
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

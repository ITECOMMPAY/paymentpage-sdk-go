package paymentpage

import (
	"encoding/json"
	"errors"
	"strconv"
)

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

type Callback struct {
	signatureHandler SignatureHandler
	callbackData     string
	parsedData       map[string]interface{}
	signature        string
}

func (c *Callback) GetPayment() interface{} {
	return c.getParamByName("payment", c.parsedData)
}

func (c *Callback) GetPaymentStatus() string {
	status := c.getParamByName("status", c.parsedData)

	return status.(string)
}

func (c *Callback) GetPaymentId() string {
	id := c.getParamByName("id", c.parsedData)

	switch id := id.(type) {
	case float64:
		return strconv.FormatFloat(id, 'f', 0, 64)
	default:
		return id.(string)
	}
}

func (c *Callback) getSignature() string {
	if c.signature == "" {
		c.signature = c.getParamByName("signature", c.parsedData).(string)
	}

	return c.signature
}

func (c *Callback) checkSignature() error {
	signature := c.getSignature()
	c.removeParam("signature", c.parsedData)

	if !c.signatureHandler.Check(signature, c.parsedData) {
		return errors.New("invalid signature")
	}

	return nil
}

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

func (c *Callback) parseCallbackData() error {
	parseError := json.Unmarshal([]byte(c.callbackData), &c.parsedData)

	if parseError != nil {
		return errors.New("invalid callback data")
	}

	return nil
}

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

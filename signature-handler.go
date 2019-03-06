package paymentpage

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"sort"
	"strconv"
	"strings"
)

// Structure for make/check signature
type SignatureHandler struct {
	// Site salt
	secret string

	// Need sort params before sign or not need
	sort bool
}

// Method for check signature
func (s *SignatureHandler) Check(signature string, params map[string]interface{}) bool {
	return signature == s.Sign(params)
}

// Setter for sort flag
func (s *SignatureHandler) SetSort(sort bool) *SignatureHandler {
	s.sort = sort

	return s
}

// Method for make signature
func (s *SignatureHandler) Sign(params map[string]interface{}) string {
	paramsToSign := s.getParamsToSign(params, "")
	arrParams := []string{}

	for _, value := range paramsToSign {
		arrParams = append(arrParams, value)
	}

	if s.sort {
		sort.Strings(arrParams)
	}

	strParams := strings.Join(arrParams, ";")
	secret := []byte(s.secret)
	message := []byte(strParams)

	hash := hmac.New(sha512.New, secret)
	hash.Write(message)

	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

// Method for preparing params
func (s *SignatureHandler) getParamsToSign(params map[string]interface{}, prefix string) map[string]string {
	paramsToSign := map[string]string{}

	for key, value := range params {
		newKey := key

		if prefix != "" {
			newKey = concat(concat(prefix, ":"), key)
		}

		preparedValue := ""

		switch value := value.(type) {
		case string:
			preparedValue = value
		case int:
			preparedValue = strconv.Itoa(value)
		case bool:
			preparedValue = getStringBool(value)
		case map[string]interface{}:
			paramsToSign = mergeMaps(paramsToSign, s.getParamsToSign(value, newKey))
		}

		if preparedValue != "" {
			paramsToSign[newKey] = concat(concat(newKey, ":"), preparedValue)
		}
	}

	return paramsToSign
}

// Constructor for SignatureHandler structure
func NewSignatureHandler(secret string) *SignatureHandler {
	signatureHandler := SignatureHandler{secret, true}

	return &signatureHandler
}

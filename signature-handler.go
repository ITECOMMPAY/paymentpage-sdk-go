package paymentpage

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"sort"
	"strconv"
	"strings"
)

type SignatureHandler struct {
	secret string
	sort   bool
}

func (s *SignatureHandler) Check(signature string, params map[string]interface{}) bool {
	return signature == s.Sign(params)
}

func (s *SignatureHandler) SetSort(sort bool) *SignatureHandler {
	s.sort = sort

	return s
}

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

func NewSignatureHandler(secret string) *SignatureHandler {
	signatureHandler := SignatureHandler{secret, true}

	return &signatureHandler
}

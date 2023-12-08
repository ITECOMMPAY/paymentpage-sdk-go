package paymentpage

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
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
	strParams := s.getStringParamsToSign(params)
	secret := []byte(s.secret)
	message := []byte(strParams)

	hash := hmac.New(sha512.New, secret)
	hash.Write(message)

	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

// Method for make string with params for signature
func (s *SignatureHandler) getStringParamsToSign(params map[string]interface{}) string {
	paramsToSign := s.getParamsToSign(params, "")
	arrParams := []string{}

	for _, value := range paramsToSign {
		arrParams = append(arrParams, value)
	}

	if s.sort {
		sort.Strings(arrParams)
	}

	return strings.Join(arrParams, ";")
}

// Method for preparing params
func (s *SignatureHandler) getParamsToSign(params map[string]interface{}, prefix string) map[string]string {
	paramsToSign := map[string]string{}
	subParamsToSign := map[string]interface{}(nil)

	for key, value := range params {
		if s.isIgnoredKey(key) {
			continue
		}

		newKey := strings.ReplaceAll(key, ":", "::")

		if prefix != "" {
			newKey = concat(concat(prefix, ":"), key)
		}

		preparedValue := ""
		valueIterable := false

		switch value := value.(type) {
		case bool:
			preparedValue = getStringBool(value)
		case int:
			preparedValue = strconv.Itoa(value)
		case float64:
			preparedValue = strconv.Itoa(int(value))
		case []interface{}:
			subParamsToSign = s.sliceToMap(value)
			valueIterable = true
		case map[string]interface{}:
			subParamsToSign = value
			valueIterable = true
		case nil:
			break
		default:
			preparedValue = fmt.Sprint(value)
		}

		if valueIterable && len(subParamsToSign) == 0 {
			continue
		} else if len(subParamsToSign) != 0 {
			paramsToSign = mergeMaps(paramsToSign, s.getParamsToSign(subParamsToSign, newKey))
			subParamsToSign = map[string]interface{}(nil)
		} else {
			paramsToSign[newKey] = concat(concat(newKey, ":"), preparedValue)
		}
	}

	return paramsToSign
}

// Method for check ignored keys
func (s *SignatureHandler) isIgnoredKey(key string) bool {
	//frame_mode
	ignoredKeys := []string{
		"frame_mode",
	}

	for i := 0; i < len(ignoredKeys); i++ {
		if ignoredKeys[i] == key {
			return true
		}
	}

	return false
}

// Method for convert slice to map
func (s *SignatureHandler) sliceToMap(slice []interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for i := 0; i < len(slice); i++ {
		result[strconv.Itoa(i)] = slice[i]
	}

	return result
}

// Constructor for SignatureHandler structure
func NewSignatureHandler(secret string) *SignatureHandler {
	signatureHandler := SignatureHandler{secret, true}

	return &signatureHandler
}

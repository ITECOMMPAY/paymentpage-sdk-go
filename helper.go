package paymentpage

import "unsafe"

type customString struct {
	str *byte
	len int
}

func concat(x, y string) string {
	length := len(x) + len(y)
	if length == 0 {
		return ""
	}
	b := make([]byte, length)
	copy(b, x)
	copy(b[len(x):], y)
	return goString(&b[0], length)
}

func goString(ptr *byte, length int) string {
	s := customString{str: ptr, len: length}
	return *(*string)(unsafe.Pointer(&s))
}

func getStringBool(value bool) string {
	if value {
		return "1"
	}
	return "0"
}

func mergeMaps(map1 map[string]string, map2 map[string]string) map[string]string {
	for key, value := range map2 {
		map1[key] = value
	}
	return map1
}

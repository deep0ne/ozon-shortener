package base63

import (
	"fmt"
	"strings"
)

const (
	alphabet       = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890_"
	alphabetLength = int64(len(alphabet))
	shortLength    = 10
)

func Encode(num int64) string {

	encoded := make([]byte, 0, shortLength)
	for num > 0 || len(encoded) < 10 {
		remainder := num % alphabetLength
		num /= alphabetLength
		encoded = append(encoded, alphabet[remainder])
	}

	// переворачиваем слайс и возвращаем строку
	for i := 0; i < len(encoded)/2; i++ {
		j := len(encoded) - i - 1
		encoded[i], encoded[j] = encoded[j], encoded[i]
	}

	return string(encoded)
}

func Decode(encoded string) (int64, error) {
	var decoded int64
	encodedLength := len(encoded)
	for i := 0; i < encodedLength; i++ {
		c := encoded[i]
		value := int64(strings.IndexByte(alphabet, c))
		if value == -1 {
			return 0, fmt.Errorf("invalid character: '%c'", c)
		}
		decoded = decoded*alphabetLength + value
	}
	return decoded, nil
}

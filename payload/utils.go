package payload

import (
	"encoding/base64"
	"math"
	"strconv"
)

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func floatToHex(x float64) string {
	var result []byte
	quotient := int(x)
	fraction := x - float64(quotient)

	for quotient > 0 {
		quotient = int(x / 16)
		remainder := int(x - (float64(quotient) * 16))

		if remainder > 9 {
			result = append([]byte{byte(remainder + 55)}, result...)
		} else {
			for _, c := range strconv.Itoa(int(remainder)) {
				result = append([]byte{byte(c)}, result...)
			}
		}

		x = float64(quotient)
	}

	if fraction == 0 {
		return string(result)
	}

	result = append(result, '.')

	for fraction > 0 {
		fraction = fraction * 16
		integer := int(fraction)
		fraction = fraction - float64(integer)

		if integer > 9 {
			result = append(result, byte(integer+55))
		} else {
			for _, c := range strconv.Itoa(int(integer)) {
				result = append(result, byte(c))
			}
		}
	}

	return string(result)
}

func a(b, c, d float64) float64 {
	return b*(d-c)/255 + c
}

func b(a int) float64 {
	if a%2 == 1 {
		return -1.0
	}
	return 0.0
}

func btoa(str []byte) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func timeToBytes(val uint32) []int {
	r := make([]int, 4)
	for i := uint32(0); i < 4; i++ {
		r[i] = int((val >> (8 * i)) & 0xff)
	}
	return r
}

func atob(input string) string {
	data, err := base64.RawStdEncoding.DecodeString(input)
	if err != nil {
		return ""
	}
	return string(data)
}

func charCodeAt(a string, i int) int {
	return int(a[i])
}

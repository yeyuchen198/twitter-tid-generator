package payload

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

const totalTime = 4096.0

// GenerateHeader - generates an x-transaction-id header
func GenerateHeader(path, method, key string, frames [][][]int) string {
	keyBytes := []int{}
	key = atob(key)
	for i := 0; i < len(key); i++ {
		keyBytes = append(keyBytes, charCodeAt(key, i))
	}
	timeNow := uint32((time.Now().UnixMilli() - 1682924400*1000) / 1000)
	// timeNow := uint32(17615204)
	timeNowBytes := timeToBytes(timeNow)

	row := frames[keyBytes[5]%4][keyBytes[2]%16]
	targetTime := float64(keyBytes[12]%16*(keyBytes[14]%16)*(keyBytes[7]%16)) / totalTime
	fromColor := []float64{float64(row[0]), float64(row[1]), float64(row[2]), 1.0}
	toColor := []float64{float64(row[3]), float64(row[4]), float64(row[5]), 1.0}
	fromRotation := []float64{0.0}
	toRotation := []float64{a(float64(row[6]), 60.0, 360.0)}
	row = row[7:]
	curves := [4]float64{}
	for i := 0; i < len(row); i++ {
		curves[i] = a(float64(row[i]), b(i), 1.0)
	}
	c := &cubic{Curves: curves}
	val := c.getValue(targetTime)
	color := interpolate(fromColor, toColor, val)
	rotation := interpolate(fromRotation, toRotation, val)
	matrix := convertRotationToMatrix(rotation[0])
	strArr := []string{}
	for i := 0; i < len(color)-1; i++ {
		strArr = append(strArr, hex.EncodeToString([]byte{byte(math.Round(color[i]))}))
	}
	for i := 0; i < len(matrix)-2; i++ {
		rounded := toFixed(matrix[i], 2)
		if rounded < 0 {
			rounded = -rounded
		}
		strArr = append(strArr, "0"+strings.ToLower(floatToHex(rounded)[1:]))
	}
	strArr = append(strArr, "0", "0")
	hash := sha256.Sum256([]byte(fmt.Sprintf(`%s!%s!%vbird%s`, method, path, timeNow, strings.Join(strArr, ""))))
	hashBytes := []int{}
	for i := 0; i < len(hash)-16; i++ {
		hashBytes = append(hashBytes, int(hash[i]))
	}
	xorByte := rand.Intn(256)
	// xorByte := 160
	bytes := []int{xorByte}
	bytes = append(bytes, keyBytes...)
	bytes = append(bytes, timeNowBytes...)
	bytes = append(bytes, hashBytes...)
	bytes = append(bytes, 1)
	out := []byte{}
	for i := 0; i < len(bytes); i++ {
		if i == 0 {
			// ! don't xor the xor byte
			out = append(out, byte(bytes[i]))
			continue
		}
		out = append(out, byte(bytes[i]^xorByte))
	}
	return strings.ReplaceAll(btoa(out), "=", "")
}

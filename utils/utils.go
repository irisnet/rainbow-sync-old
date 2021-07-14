package utils

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/irisnet/rainbow-sync/lib/logger"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func BuildHex(bytes []byte) string {
	return strings.ToUpper(hex.EncodeToString(bytes))
}

func ConvertErr(height int64, txHash, errTag string, err error) error {
	return fmt.Errorf("%v-%v-%v-%v", err.Error(), errTag, height, txHash)
}

func GetErrTag(err error) string {
	slice := strings.Split(err.Error(), "-")
	if len(slice) == 4 {
		return slice[1]
	}
	return ""
}

func Min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func ParseFloat(s string, bit ...int) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		logger.Error("common.ParseFloat error", logger.String("value", s))
		return 0
	}

	if len(bit) > 0 {
		return RoundFloat(f, bit[0])
	}
	return f
}

func RoundFloat(num float64, bit int) (i float64) {
	format := "%" + fmt.Sprintf("0.%d", bit) + "f"
	s := fmt.Sprintf(format, num)
	i, err := strconv.ParseFloat(s, 0)
	if err != nil {
		logger.Error("common.RoundFloat error", logger.String("format", format))
		return 0
	}
	return i
}

// Intn returns, as an int, a non-negative pseudo-random number in [0,n)
// from the default Source.
// It panics if n <= 0.
func RandInt(n int) int {
	rand.NewSource(time.Now().Unix())
	return rand.Intn(n)
}

func MarshalJsonIgnoreErr(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

func UnMarshalJsonIgnoreErr(data string, v interface{}) {
	json.Unmarshal([]byte(data), &v)
}

func RemoveDuplicatesFromSlice(data []string) (result []string) {
	tempSet := make(map[string]string, len(data))
	for _, val := range data {
		if _, ok := tempSet[val]; ok || val == "" {
			continue
		}
		tempSet[val] = val
	}
	for one := range tempSet {
		result = append(result, one)
	}
	return
}

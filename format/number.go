package format

import (
	"fmt"
	"strconv"
)

// ToFixedWithTwoDigits 保留两位小数
func ToFixedWithTwoDigits(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

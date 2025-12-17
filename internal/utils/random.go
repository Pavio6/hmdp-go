package utils

import (
	"crypto/rand"
	"math/big"
)

// RandomString 函数生成指定长度的随机字符串。
// 字符包括小写字母和数字。
// 它使用 crypto/rand 函数以获得更好的随机性。
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			result[i] = charset[0]
			continue
		}
		result[i] = charset[num.Int64()]
	}
	return string(result)
}

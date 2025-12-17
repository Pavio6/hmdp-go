package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const verifyCodeMax = 1000000

// GenerateVerifyCode 返回一个以零填充的 6 位数字验证码
func GenerateVerifyCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(verifyCodeMax))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

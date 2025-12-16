package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/google/uuid"
)

// Encode hashes a password with a random salt.
func Encode(password string) string {
	salt := uuid.NewString()
	return encodeWithSalt(password, salt)
}

// Matches verifies a password against the encoded representation.
func Matches(encodedPassword, rawPassword string) (bool, error) {
	if encodedPassword == "" || rawPassword == "" {
		return false, nil
	}
	parts := strings.Split(encodedPassword, "@")
	if len(parts) != 2 {
		return false, errors.New("invalid password format")
	}
	salt := parts[0]
	return encodeWithSalt(rawPassword, salt) == encodedPassword, nil
}

func encodeWithSalt(password, salt string) string {
	sum := md5.Sum([]byte(password + salt))
	return salt + "@" + hex.EncodeToString(sum[:])
}

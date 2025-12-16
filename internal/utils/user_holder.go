package utils

import (
	"sync/atomic"

	"hmdp-backend/internal/dto"
)

var userHolder atomic.Value

// SaveUser replicates the Java ThreadLocal behaviour (process-wide placeholder).
func SaveUser(user *dto.UserDTO) {
	userHolder.Store(user)
}

// GetUser returns the stored user, if any.
func GetUser() *dto.UserDTO {
	if v := userHolder.Load(); v != nil {
		if user, ok := v.(*dto.UserDTO); ok {
			return user
		}
	}
	return nil
}

// RemoveUser clears the stored user.
func RemoveUser() {
	userHolder.Store((*dto.UserDTO)(nil))
}

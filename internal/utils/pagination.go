package utils

import "strconv"

func ParsePage(value string, defaultVal int) int {
	if value == "" {
		return defaultVal
	}
	if v, err := strconv.Atoi(value); err == nil && v > 0 {
		return v
	}
	return defaultVal
}

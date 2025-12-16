package utils

import "strconv"

// ParsePage safely converts query values to positive page numbers.
func ParsePage(value string, defaultVal int) int {
    if value == "" {
        return defaultVal
    }
    if v, err := strconv.Atoi(value); err == nil && v > 0 {
        return v
    }
    return defaultVal
}

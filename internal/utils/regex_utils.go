package utils

import "regexp"

// IsPhoneInvalid replicates RegexUtils#isPhoneInvalid.
func IsPhoneInvalid(phone string) bool {
	return mismatch(phone, PHONE_REGEX)
}

// IsEmailInvalid replicates RegexUtils#isEmailInvalid.
func IsEmailInvalid(email string) bool {
	return mismatch(email, EMAIL_REGEX)
}

// IsCodeInvalid replicates RegexUtils#isCodeInvalid.
func IsCodeInvalid(code string) bool {
	return mismatch(code, VERIFY_CODE_REGEX)
}

func mismatch(value, pattern string) bool {
	if value == "" {
		return true
	}
	matched, _ := regexp.MatchString(pattern, value)
	return !matched
}

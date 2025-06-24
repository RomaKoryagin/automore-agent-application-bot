package helpers

import "regexp"

func IsPhoneValid(phone string) bool {
	pattern := regexp.MustCompile(`^7\d{10}$`)
	return pattern.MatchString(phone)
}

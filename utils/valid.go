package untils

import (
	"net/mail"
	"regexp"
)

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsValidUsername(username string) bool {
	if len(username) < 1 || len(username) > 30 {
		return false
	}
	match, _ := regexp.MatchString("^[a-zA-Z0-9](?:[a-zA-Z0-9_]*[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9]+)*$", username)
	return match
}
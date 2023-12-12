package untils

import (
	"regexp"
)

func IsValidEmail(email string) bool {
	
	match, _ := regexp.MatchString("^[a-z0-9._%+-]+@[a-z0-9.-]+.[a-z]$", email)
	return match
}

func IsValidUsername(username string) bool {
	if len(username) < 1 || len(username) > 30 {
		return false
	}
	match, _ := regexp.MatchString("^[a-zA-Z0-9](?:[a-zA-Z0-9_]*[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9]+)*$", username)
	return match
}
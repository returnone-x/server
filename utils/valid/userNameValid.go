package valid

import "regexp"

func IsValidUsername(username string) bool {
	if len(username) < 1 || len(username) > 30 {
		return false
	}
	match, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	return match
}

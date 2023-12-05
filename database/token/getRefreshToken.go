package tokenDatabase

import (
	"returnone/config"
)

func GetRefreshToken(
	token string,
) int {
	var count int

	sqlString := `SELECT COUNT(*) FROM tokens WHERE token = $1;`
	config.DB.QueryRow(sqlString, token).Scan(&count)

	return count
}

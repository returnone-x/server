package tokenDatabase

import (
	db "returnone/config"
	"returnone/models/key"
)

func GetTokenData(id string) (token_data keyModles.TokenType, err error) {
	
	sqlString := `SELECT id, used_time, user_agent, ip, create_at, update_at FROM tokens WHERE id = $1;`
	err = db.DB.QueryRow(sqlString, id).Scan(
		&token_data.Id,
		&token_data.Used_time,
		&token_data.User_agent,
		&token_data.Ip,
		&token_data.Create_at,
		&token_data.Update_at,)
		
	return token_data, err
}


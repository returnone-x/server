package auth

import (
	"log"
	tokenDatabase "github.com/returnone-x/server/database/tokens"
	utils "github.com/returnone-x/server/utils"

	"github.com/gofiber/fiber/v2"
)

func SetLoginCookies(user_id string, c *fiber.Ctx) (access_token string, refresh_token string, error_message string, err error,) {
	// btw actually access_token_id not been used in this program (i dont know if i will use it in the future)
	access_token_id := utils.GenerateTokenId()
	// generate token id for check the token have been used or not
	refresh_token_id := utils.GenerateTokenId()

	//generate Jwt token (exp set on the handler)
	access_token, access_token_err := utils.GenerateJwtToken(user_id, access_token_id, 0,"accessToken", access_token_exp.Unix())
	refresh_token, refresh_token_err := utils.GenerateJwtToken(user_id, refresh_token_id, 1,"refreshToken", refresh_token_exp.Unix())
	
	//handle errors
	if refresh_token_err != nil {
		log.Println("| Path:", c.Path(), "| Data:", user_id, "| Message:", refresh_token_err)
		return "", "", "Error generating refresh token", refresh_token_err
	}
	if access_token_err != nil {
		log.Println("| Path:", c.Path(), "| Data:", user_id, "| Message:", access_token_err)
		return "", "", "Error generating access token", access_token_err
	}

	ip := c.IP()
	user_agent := c.Get("User-Agent")

	result, err := tokenDatabase.CreateToken(refresh_token_id, user_agent, ip)
	count, _ := result.RowsAffected()
	
	if err != nil || count == 0 {
		log.Println("| Path:", c.Path(), "| Data:", user_id, "| Message:", err)
		return "", "", "Error saving refresh token", err
	}

	return access_token, refresh_token, "", nil
}

func SetRefreshCookies(user_id string, refresh_token_id string, used_time int, c *fiber.Ctx) (access_token string, refresh_token string, error_message string, err error,) {
	// btw actually access_token_id not been used in this program (i dont know if i will use it in the future)
	access_token_id := utils.GenerateTokenId()

	//generate Jwt token (exp set on the handler)
	access_token, access_token_err := utils.GenerateJwtToken(user_id, access_token_id, 0, "accessToken", access_token_exp.Unix())
	refresh_token, refresh_token_err := utils.GenerateJwtToken(user_id, refresh_token_id, used_time, "refreshToken", refresh_token_exp.Unix())

	//handle errors
	if refresh_token_err != nil {
		log.Println("| Path:", c.Path(), "| Data:", user_id, "| Message:", refresh_token_err)
		return "", "", "Error generating refresh token", refresh_token_err
	}
	if access_token_err != nil {
		log.Println("| Path:", c.Path(), "| Data:", user_id, "| Message:", access_token_err)
		return "", "", "Error generating access token", access_token_err
	}

	result, err := tokenDatabase.UpdateToken(refresh_token_id, used_time)
	count, _ := result.RowsAffected()

	if err != nil || count == 0 {
		log.Println("| Path:", c.Path(), "| Data:", user_id, "| Message:", err)
		return "", "", "Error saving refresh token", err
	}
	
	return access_token, refresh_token, "", nil
}
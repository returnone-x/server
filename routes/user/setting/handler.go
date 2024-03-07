package setting

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	userDatabase "github.com/returnone-x/server/database/user"
	userSettingDatabase "github.com/returnone-x/server/database/user/setting"
	utils "github.com/returnone-x/server/utils"
)

func ResetPassword(c *fiber.Ctx) error {
	var data map[string]string

	err := c.BodyParser(&data)

	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Invalid post request",
			})
	}

	if data["new_password"] == "" {
		return c.Status(400).JSON(utils.RequestValueValid("new_password"))
	}

	// get user_id from accessToken cookie
	token := c.Locals("access_token_context").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	user_id := claims["user_id"].(string)

	user_data, err := userDatabase.GetUserDetil(user_id)

	if err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(utils.ErrorMessage("Error get user detil", err))
	}

	if user_data.Password == "" {
		hash_password, err := utils.HashPassword(data["new_password"])

		if err != nil {
			return c.Status(500).JSON(utils.ErrorMessage("Error hash user password", err))
		}

		result, err := userSettingDatabase.UpdateUserPassword(user_id, hash_password)

		if err != nil {
			return c.Status(500).JSON(utils.ErrorMessage("Error update password", err))
		}

		// check does it really update
		row_affected, _ := result.RowsAffected()

		// if not update return server error or return 200 status code
		if row_affected == 0 {
			return c.Status(500).JSON(utils.ErrorMessage("Error update password", err))
		}

		return c.Status(200).JSON(fiber.Map{
			"status":  "successful",
			"message": "successful create your password",
		})
	}

	if data["old_password"] == "" {
		return c.Status(400).JSON(utils.RequestValueValid("old password"))
	}

	verify_password := utils.CheckPasswordHash(data["old_password"], user_data.Password)
	if !verify_password {
		return c.Status(401).JSON(utils.RequestValueValid("old password"))
	}

	hash_password, err := utils.HashPassword(data["new_password"])
	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error hash new user password", err))
	}

	result, err := userSettingDatabase.UpdateUserPassword(user_id, hash_password)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error update password", err))
	}
	// check does it really update
	row_affected, _ := result.RowsAffected()
	// if not update return server error or return 200 status code
	if row_affected == 0 {
		return c.Status(500).JSON(utils.ErrorMessage("Error update password", err))
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "successful",
		"message": "successful update password",
	})
}

func ResetAvatar(c *fiber.Ctx) error {

	var data map[string]string

	err := c.BodyParser(&data)

	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Invalid post request",
			})
	}

	if data["avatar"] == "" {
		return c.Status(400).JSON(utils.RequestValueValid("avatar"))
	}

	if len(data["avatar"]) > 255 {
		return c.Status(400).JSON(utils.RequestValueValid("avatar"))
	}

	token := c.Locals("access_token_context").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	user_id := claims["user_id"].(string)

	result, err := userSettingDatabase.UpdateUserAvatar(user_id, data["avatar"])
	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error update avatar", err))
	}
	// check does it really update
	row_affected, _ := result.RowsAffected()
	// if not update return server error or return 200 status code
	if row_affected == 0 {
		return c.Status(500).JSON(utils.ErrorMessage("Error update avatar", err))
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "successful",
		"message": "successful update user avatar",
	})
}

func ResetUsername(c *fiber.Ctx) error {

	var data map[string]string

	err := c.BodyParser(&data)

	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Invalid post request",
			})
	}

	if !utils.IsValidUsername(data["username"]) {
		return c.Status(400).JSON(utils.RequestValueValid("new username"))
	}

	// check the user name has already been used
	if userDatabase.CheckUserNameExist(data["username"]) != 0 {
		return c.Status(400).JSON(utils.RequestValueInUse("new username"))
	}

	token := c.Locals("access_token_context").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	user_id := claims["user_id"].(string)

	result, err := userSettingDatabase.UpdateUsername(user_id, data["username"])
	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error update username", err))
	}
	// check does it really update
	row_affected, _ := result.RowsAffected()
	// if not update return server error or return 200 status code
	if row_affected == 0 {
		return c.Status(500).JSON(utils.ErrorMessage("Error update username", err))
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "successful",
		"message": "successful update username",
	})
}

func ResetDisplayName(c *fiber.Ctx) error {

	var data map[string]string

	err := c.BodyParser(&data)

	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Invalid post request",
			})
	}

	if data["display_name"] == "" {
		return c.Status(400).JSON(utils.RequestValueValid("display name"))
	}

	if len(data["display_name"]) > 30 {
		return c.Status(400).JSON(utils.RequestValueValid("display name"))
	}

	token := c.Locals("access_token_context").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	user_id := claims["user_id"].(string)

	result, err := userSettingDatabase.UpdateDisplayName(user_id, data["display_name"])
	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error update display name", err))
	}
	// check does it really update
	row_affected, _ := result.RowsAffected()
	// if not update return server error or return 200 status code
	if row_affected == 0 {
		return c.Status(500).JSON(utils.ErrorMessage("Error update display name", err))
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "successful",
		"message": "successful update display name",
	})
}

func ResetTotp(c *fiber.Ctx) error {

	var data map[string]string

	err := c.BodyParser(&data)

	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Invalid post request",
			})
	}

	if data["totp"] == "" {
		return c.Status(400).JSON(utils.RequestValueValid("totp"))
	}

	if len(data["totp"]) > 30 {
		return c.Status(400).JSON(utils.RequestValueValid("totp"))
	}

	token := c.Locals("access_token_context").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	user_id := claims["user_id"].(string)

	result, err := userSettingDatabase.UpdateTotp(user_id, data["totp"])
	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error update totp", err))
	}
	// check does it really update
	row_affected, _ := result.RowsAffected()
	// if not update return server error or return 200 status code
	if row_affected == 0 {
		return c.Status(500).JSON(utils.ErrorMessage("Error update totp", err))
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "successful",
		"message": "successful update totp",
	})
}

func CheckPassword(password string) bool {
	if password == "" {
		return false
	} else {
		return true
	}
}

func GetUser(c *fiber.Ctx) error {

	token := c.Locals("access_token_context").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	user_id := claims["user_id"].(string)

	user_data, err := userDatabase.GetUserDetil(user_id)

	if err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(utils.ErrorMessage("Error get user detil", err))
	}

	user_data_detil := map[string]interface{}{
		"id":             user_data.Id,
		"email":          user_data.Email,
		"phone":          user_data.Phone_country,
		"phone_country":  user_data.Phone_country,
		"password":       CheckPassword(user_data.Password),
		"email_verify":   user_data.Email_verify,
		"phone_verify":   user_data.Phone_verify,
		"avatar":         user_data.Avatar,
		"display_name":   user_data.Display_name,
		"username":       user_data.Username,
		"github_connect": user_data.Github_connect,
		"google_connect": user_data.Google_connect,
		"email_2fa":      user_data.Email_2fa,
		"phone_2fa":      user_data.Phone_2fa,
		"totp_2fa":       user_data.Totp_2fa,
		"default_2fa":    user_data.Default_2fa,
		"bio": user_data.Bio,
		"public_email": user_data.Public_email,
		"pronouns": user_data.Pronouns,
		"related_links": user_data.Related_links,
		"create_at":      user_data.Create_at,
		"update_at":      user_data.Update_at,
	}
	return c.Status(200).JSON(fiber.Map{
		"data":    user_data_detil,
		"status":  "successful",
		"message": "successful get user setting detil",
	})
}

func ResetBio(c *fiber.Ctx) error {

	var data map[string]string

	err := c.BodyParser(&data)

	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Invalid post request",
			})
	}

	if data["bio"] == "" {
		return c.Status(400).JSON(utils.RequestValueValid("bio"))
	}

	if len(data["bio"]) > 150 {
		return c.Status(400).JSON(utils.RequestValueValid("bio"))
	}

	token := c.Locals("access_token_context").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	user_id := claims["user_id"].(string)

	result, err := userSettingDatabase.UpdateBio(user_id, data["bio"])
	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error update bio", err))
	}
	// check does it really update
	row_affected, _ := result.RowsAffected()
	// if not update return server error or return 200 status code
	if row_affected == 0 {
		return c.Status(500).JSON(utils.ErrorMessage("Error update bio", err))
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "successful",
		"message": "successful update bio",
	})
}

func ResetPublicEmail(c *fiber.Ctx) error {

	var data map[string]string

	err := c.BodyParser(&data)

	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Invalid post request",
			})
	}

	if data["public_email"] == "" {
		return c.Status(400).JSON(utils.RequestValueValid("public_email"))
	}

	if len(data["public_email"]) > 150 {
		return c.Status(400).JSON(utils.RequestValueValid("public_email"))
	}

	token := c.Locals("access_token_context").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	user_id := claims["user_id"].(string)

	result, err := userSettingDatabase.UpdatePublicEmail(user_id, data["public_email"])
	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error update public email", err))
	}
	// check does it really update
	row_affected, _ := result.RowsAffected()
	// if not update return server error or return 200 status code
	if row_affected == 0 {
		return c.Status(500).JSON(utils.ErrorMessage("Error update public email", err))
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "successful",
		"message": "successful update public email",
	})
}

func ResetPronouns(c *fiber.Ctx) error {

	var data map[string]string

	err := c.BodyParser(&data)

	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Invalid post request",
			})
	}

	if data["pronouns"] == "" {
		return c.Status(400).JSON(utils.RequestValueValid("pronouns"))
	}

	if len(data["pronouns"]) > 20 {
		return c.Status(400).JSON(utils.RequestValueValid("pronouns"))
	}

	token := c.Locals("access_token_context").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	user_id := claims["user_id"].(string)

	result, err := userSettingDatabase.UpdatePronouns(user_id, data["pronouns"])
	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error update pronouns", err))
	}
	// check does it really update
	row_affected, _ := result.RowsAffected()
	// if not update return server error or return 200 status code
	if row_affected == 0 {
		return c.Status(500).JSON(utils.ErrorMessage("Error update pronouns", err))
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "successful",
		"message": "successful update pronouns",
	})
}

type RelateLinksRequestBody struct {
	Related_links []string `json:"related_links"`
}

func ResetRelatedLinks(c *fiber.Ctx) error {

	var data RelateLinksRequestBody

	err := c.BodyParser(&data)

	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Invalid post request",
			})
	}

	if len(data.Related_links) == 0 {
		return c.Status(400).JSON(utils.RequestValueValid("related_links"))
	}

	if len(data.Related_links) > 6 {
		return c.Status(400).JSON(utils.RequestValueValid("related_links"))
	}

	token := c.Locals("access_token_context").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	user_id := claims["user_id"].(string)

	result, err := userSettingDatabase.UpdateRelateLinks(user_id, data.Related_links)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error update related_links", err))
	}
	// check does it really update
	row_affected, _ := result.RowsAffected()
	// if not update return server error or return 200 status code
	if row_affected == 0 {
		return c.Status(500).JSON(utils.ErrorMessage("Error update related_links", err))
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "successful",
		"message": "successful update relate links",
	})
}

func ResetAllName(c *fiber.Ctx) error {

	var data map[string]string

	err := c.BodyParser(&data)

	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Invalid post request",
			})
	}

	if data["display_name"] == "" {
		return c.Status(400).JSON(utils.RequestValueValid("display_name"))
	}

	if len(data["display_name"]) > 150 {
		return c.Status(400).JSON(utils.RequestValueValid("display_name"))
	}

	if data["username"] == "" {
		return c.Status(400).JSON(utils.RequestValueValid("username"))
	}

	if len(data["username"]) > 150 || !utils.IsValidUsername(data["username"]) {
		return c.Status(400).JSON(utils.RequestValueValid("username"))
	}

	if userDatabase.CheckUserNameExist(data["username"]) != 0 {
		return c.Status(400).JSON(utils.RequestValueInUse("username"))
	}

	token := c.Locals("access_token_context").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	user_id := claims["user_id"].(string)

	result, err := userSettingDatabase.UpdateALlName(user_id, data["display_name"], data["username"])
	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error update name", err))
	}
	// check does it really update
	row_affected, _ := result.RowsAffected()
	// if not update return server error or return 200 status code
	if row_affected == 0 {
		return c.Status(500).JSON(utils.ErrorMessage("Error update name", err))
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "successful",
		"message": "successful update name",
	})
}

type ResetAllProfileRequestBody struct {
	Bio           string   `json:"bio"`
	Public_email  string   `json:"public_email"`
	Pronouns      string   `json:"pronouns"`
	Related_links []string `json:"related_links"`
}

func ResetAllProfile(c *fiber.Ctx) error {

	var data ResetAllProfileRequestBody

	err := c.BodyParser(&data)

	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Invalid post request",
			})
	}

	if len(data.Related_links) == 0 {
		return c.Status(400).JSON(utils.RequestValueValid("related_links"))
	}

	if len(data.Related_links) > 6 {
		return c.Status(400).JSON(utils.RequestValueValid("related_links"))
	}

	if len(data.Pronouns) == 0 {
		return c.Status(400).JSON(utils.RequestValueValid("pronouns"))
	}

	if len(data.Pronouns) > 20 {
		return c.Status(400).JSON(utils.RequestValueValid("pronouns"))
	}

	if len(data.Bio) == 0 {
		return c.Status(400).JSON(utils.RequestValueValid("bio"))
	}

	if len(data.Bio) > 150 {
		return c.Status(400).JSON(utils.RequestValueValid("bio"))
	}

	if len(data.Public_email) == 0 {
		return c.Status(400).JSON(utils.RequestValueValid("public_email"))
	}

	if len(data.Public_email) > 150 || !utils.IsValidEmail(data.Public_email) {
		return c.Status(400).JSON(utils.RequestValueValid("public_email"))
	}

	token := c.Locals("access_token_context").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	user_id := claims["user_id"].(string)

	result, err := userSettingDatabase.UpdateAllProfile(user_id, data.Bio, data.Public_email, data.Pronouns, data.Related_links)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error update profile", err))
	}
	// check does it really update
	row_affected, _ := result.RowsAffected()
	// if not update return server error or return 200 status code
	if row_affected == 0 {
		return c.Status(500).JSON(utils.ErrorMessage("Error update profile", err))
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "successful",
		"message": "successful update profile",
	})
}

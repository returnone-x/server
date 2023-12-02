package userModles

type User2fa struct {
	User_id     string `json:"user_id"`
	Uhone_2fa   bool   `json:"phone_2fa"`
	Email_2fa   bool   `json:"email_2fa"`
	Totp_2fa    bool   `json:"totp_2fa"`
	Totp        string `json:"totp"`
	Default_2fa int    `json:"default_2fa"`
}

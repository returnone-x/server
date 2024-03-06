package setting

import (
	"github.com/gofiber/fiber/v2"
)

func Setup(app fiber.Router) {
	app.Post("/reset/password", ResetPassword)
	app.Post("/reset/avatar", ResetAvatar)
	app.Post("/reset/username", ResetUsername)
	app.Post("/reset/displayname", ResetDisplayName)
	app.Post("/reset/totp", ResetTotp)
	app.Post("/reset/publicEmail", ResetPublicEmail)
	app.Post("/reset/pronouns", ResetPronouns)
	app.Post("/reset/bio", ResetBio)
	app.Post("/reset/relatedlinks", ResetRelatedLinks)

	app.Post("/resetall/name", ResetAllName)
	app.Post("/resetall/profile", ResetAllProfile)


	app.Get("/detil", GetUser)
}

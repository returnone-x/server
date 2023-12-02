package crypto

import (
	"os"
	"time"
	"github.com/golang-jwt/jwt"
)

func GenerateJwtToken(user_id string, Audience string,exp_time int64) (string, error) {

	SecretKey := os.Getenv("JWT_SECRET")

    claims := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.StandardClaims{
        Id:user_id,
		Audience: Audience,
		IssuedAt: time.Now().Unix(),
        ExpiresAt: exp_time,
	})

    token, err := claims.SignedString([]byte(SecretKey))

	return token, err
}


package untils

import (
	"os"
	"time"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJwtToken(user_id string, Audience string, Subject string, exp_time int64) (string, error) {

	SecretKey := os.Getenv("JWT_SECRET")

	claims := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.StandardClaims{
		Id:        user_id,
		Audience:  Audience,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: exp_time,
		Subject: Subject,
	})

	token, err := claims.SignedString([]byte(SecretKey))

	return token, err
}

func VerifyJwtToken(jwt_token *jwt.Token, id string) (bool, error) {

	claims := jwt_token.Claims.(jwt.MapClaims)
	jwt_token_user_id := claims["user_id"].(string)

	return jwt_token_user_id == id, nil
}
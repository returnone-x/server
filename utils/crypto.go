package untils

import (
	"os"
	"strconv"
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

func GenerateJwtToken(user_id string, token_id string, used_time int, subject string, exp_time int64) (string, error) {

	SecretKey := os.Getenv("JWT_SECRET")

	token := jwt.New(jwt.SigningMethodHS512)
	
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user_id
	claims["used_time"] = strconv.Itoa(used_time)
	claims["token_id"] = token_id
	claims["subject"] = subject
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	
	t, err := token.SignedString([]byte(SecretKey))

	return t, err
}

func VerifyJwtToken(jwt_token *jwt.Token, id string) (bool, error) {

	claims := jwt_token.Claims.(jwt.MapClaims)
	jwt_token_user_id := claims["user_id"].(string)

	return jwt_token_user_id == id, nil
}
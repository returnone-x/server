package untils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"time"
)

func GenerateUserAccountId() string {
	snowflake.Epoch = time.Date(2023, 12, 2, 10, 31, 32, 0, time.UTC).UnixMilli()
	node, _ := snowflake.NewNode(1)
	id := node.Generate()

	return fmt.Sprint(id)
}

func GenerateTokenId() string {
	snowflake.Epoch = time.Date(2023, 12, 2, 10, 31, 32, 0, time.UTC).UnixMilli()
	node, _ := snowflake.NewNode(2)
	id := node.Generate()

	return fmt.Sprint(id)
}

func GenerateRandomBase64String() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	state := base64.URLEncoding.EncodeToString(randomBytes)
	return state, nil
}

package Generate

import (
	"fmt"
	"time"
	"github.com/bwmarrin/snowflake"
)

func GenerateUserAccountId() string {
	snowflake.Epoch = time.Date(2023, 12, 2, 10, 31, 32, 0, time.UTC).UnixMilli()
	node, _ := snowflake.NewNode(1)
	id := node.Generate()

	return fmt.Sprint(id)
}
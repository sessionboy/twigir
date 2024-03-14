package utils

import (
	"github.com/bwmarrin/snowflake"
)

func GenerateId() int64 {
	node, _ := snowflake.NewNode(1)
	id := node.Generate().Int64()
	return id
}

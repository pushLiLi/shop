package utils

import (
	sf "github.com/bwmarrin/snowflake"
)

var node *sf.Node

func InitSnowflake(nodeID int64) {
	n, err := sf.NewNode(nodeID)
	if err != nil {
		panic(err)
	}
	node = n
}

func GenerateOrderNo() string {
	return node.Generate().String()
}

package gid

import (
	"log"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func SetNodeId(n int64) {
	var err error
	node, err = snowflake.NewNode(n)
	if err != nil {
		log.Fatalln(err)
	}

	// First test generation
	log.Printf("Node: %d, First test generation %s.\n", n, node.Generate().String())
}

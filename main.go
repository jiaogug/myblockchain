package main

import (
	"fmt"
	"myblockchain/network"
	"os"
)

func main() {

	if os.Args[1] == "runnode" {
		fmt.Println("启动节点", os.Args[2])
		nodeID := os.Args[2]
		server := network.NewServer(nodeID)

		server.Start()
	} else {
		//启动客户端
		bc := NewBlockchain()
		defer bc.db.Close()
		if os.Args[1] == "client" {
			fmt.Printf("启动客户端\n")
		}

		cli := CLI{bc}
		cli.Run()
	}
}

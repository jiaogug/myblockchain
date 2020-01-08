package main

import (
	//powr "myblockchain/pow"
	//"strconv"

	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"unsafe"
)

type client chan<- string

var nodenum = 0
var leaving = make(chan client)

// CLI responsible for processing command line arguments
type CLI struct {
	bc *Blockchain
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  addblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println("  printchain - print all the blocks of the blockchain")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) addBlock(data string) {
	cli.bc.AddBlock(data)
	fmt.Println("Success!")
}

func (cli *CLI) printChain() {
	bci := cli.bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		//pow := powr.NewProofOfWork(block)
		// pbft := pbfter.NewServer("Apple")
		//fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

// Run parses command line arguments and processes commands
func (cli *CLI) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block data")

	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	default:
		cli.printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func sendreq() {
	postmsg := make(map[string]interface{})
	postmsg["clientID"] = "ahuhwi"
	postmsg["operation"] = "GetMyname"
	postmsg["timestamp"] = 859381532

	bytesData, err := json.Marshal(postmsg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	reader := bytes.NewReader(bytesData)
	url := "http://localhost:1111/req"
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println(*str)
}

func handlereq() bool {
	listener, err := net.Listen("tcp", "localhost:8234")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		// 阻塞等待客户端连接
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		//增加超时，如果长时间没链接，认为共识未达成‘return false
		nodenum++
		fmt.Println(" There be a cli connected ...")
		//客户端只有收到3条回复就可以判定共识达成
		if nodenum > 2 {
			return true
		}
		go handleConn(conn)
	}

}

func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)
	//go clientWriter(conn, ch)
	who := conn.RemoteAddr().String() // ip地址转字符串
	ch <- "you are " + who

	leaving <- ch //转到33行

	conn.Close()
}

// func clientWriter(conn net.Conn, ch <-chan string) {
// 	for msg := range ch {
// 		fmt.Println(msg)
// 		fmt.Fprintln(conn, msg)
// 	}
// }

// func handleClient(conn net.Conn) {
// 	defer conn.Close()
// }
// server := ":8234"
// tcpAddr, err := net.ResolveTCPAddr("tcp", server)
// if err != nil {
// 	log.Fatal(err)
// }

// listener, err := net.ListenTCP("tcp", tcpAddr)
// if err != nil {
// 	log.Fatal(err)
// }

// conn, err := listener.Accept() //socket 阻塞
// if err != nil {
// 	fmt.Println("client connect failed!")
// }
// fmt.Println("client connect success!")
// conn.Close()
// return true
// //handleClient(conn)

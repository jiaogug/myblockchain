package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"myblockchain/consensus"
	"net"
	"net/http"
)

var flag bool

type Server struct {
	url  string
	node *Node
}

func NewServer(nodeID string) *Server {
	node := NewNode(nodeID)
	server := &Server{node.NodeTable[nodeID], node}

	server.setRoute()
	return server
}

func (server *Server) Start() {
	//msg = "{\"clientID\":\"ahnhwi\",\"operation\":\"GetMyname\",\"timestamp\":859381532}"
	fmt.Printf("Server will be started at %s...\n", server.url)
	if err := http.ListenAndServe(server.url, nil); err != nil {
		fmt.Println(err)
		return
	}

}

func (server *Server) setRoute() {
	http.HandleFunc("/req", server.getReq)
	http.HandleFunc("/preprepare", server.getPrePrepare)
	http.HandleFunc("/prepare", server.getPrepare)
	http.HandleFunc("/commit", server.getCommit)
	//http.HandleFunc("/reply", server.getReply)
}

// func (server *Server) sendreq(msg string) {
// 	var mesg consensus.RequestMsg
// 	msgdata := strings.NewReader(msg)
// 	err := json.NewDecoder(msgdata).Decode(&mesg) // string 转json 解码为mesg 将消息发出即可
// 	if err != nil {
// 		fmt.Println("err msg")
// 		fmt.Println(err)
// 		return
// 	}

// 	fmt.Println(mesg)                //{859381532 ahnhwi GetMyname 0}
// 	server.node.MsgEntrance <- &mesg // 把信息丢到通道中，后面是pbft流程
// }

func (server *Server) getReq(writer http.ResponseWriter, request *http.Request) {
	var msg consensus.RequestMsg
	err := json.NewDecoder(request.Body).Decode(&msg)
	if err != nil {
		fmt.Println("err msg")
		fmt.Println(err)
		return
	}
	fmt.Println("------------req msg-----------")
	fmt.Println(msg) //{859381532 ahnhwi GetMyname 0}
	fmt.Println("------------req-----------")
	server.node.MsgEntrance <- &msg // 把信息丢到通道中
}

func (server *Server) getPrePrepare(writer http.ResponseWriter, request *http.Request) {
	var msg consensus.PrePrepareMsg
	err := json.NewDecoder(request.Body).Decode(&msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	server.node.MsgEntrance <- &msg
}

func (server *Server) getPrepare(writer http.ResponseWriter, request *http.Request) {
	var msg consensus.VoteMsg
	err := json.NewDecoder(request.Body).Decode(&msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	server.node.MsgEntrance <- &msg
}

func (server *Server) getCommit(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("收到commit并发出")
	var msg consensus.VoteMsg
	err := json.NewDecoder(request.Body).Decode(&msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	server.node.MsgEntrance <- &msg
}

func (server *Server) getReply(writer http.ResponseWriter, request *http.Request) {
	//fmt.Println("----------	(server *Server) getReply-----------")
	var msg consensus.ReplyMsg
	err := json.NewDecoder(request.Body).Decode(&msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	server.node.GetReply(&msg)
}

func send(url string, msg []byte) {
	buff := bytes.NewBuffer(msg)
	//fmt.Println("-----------", buff, "-----------")
	http.Post("http://"+url, "application/json", buff)
}
func sendmsgtoclient(url string, msg []byte) {
	conn, err := net.Dial("tcp", url)

	if err != nil {
		//log.Fatal(err)
		fmt.Println("warining, connectfailed2")
		//conn.Close()
		return
	}
	conn.Close()
}

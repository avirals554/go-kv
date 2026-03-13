package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

var port = ""
var active, err = os.OpenFile("active_node.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)

func load_port() {
	active.Seek(0, 0)
	scanner := bufio.NewScanner(active)
	for scanner.Scan() {
		line := scanner.Text()
		port = line
	}
}
func main() {

	listen_client, err := net.Listen("tcp", "localhost:5378")
	if err != nil {
		fmt.Println("client side problem ")
	}
	for {
		clientconn, err := listen_client.Accept()
		if err != nil {
			fmt.Println("there is something wrong with the client conn ")
			return
		}
		load_port()
		leaderconn, err := net.Dial("tcp", "localhost:"+port)
		if err != nil {
			fmt.Println("something is wrong  with the connection to the leader ")
			return
		}
		if err == nil {
			fmt.Println("i am trying to connect to the leader ")
		}
		go io.Copy(clientconn, leaderconn) // server → client
		go io.Copy(leaderconn, clientconn) // client → server
	}
}

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

var mu sync.Mutex
var store = make(map[string]string)
var db, err = os.OpenFile("val.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)

func load_from_disk() {
	scanner := bufio.NewScanner(db)
	for scanner.Scan() {
		line := scanner.Text()
		command := strings.Split(line, " ")
		if command[0] == "SET" {
			mu.Lock()
			store[command[1]] = command[2]
			mu.Unlock()
		}
	}
}
func makeconnection(conn net.Conn) {
	fmt.Println("we are getting a connection from  ", conn.RemoteAddr())
	conn.Write([]byte("tcp connection has started ......."))

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("something went went wrong conn.read()")
			break
		}
		message := strings.Split(string(buf[:n]), " ")
		if message[0] == "GET" {
			mu.Lock()
			value, ok := store[message[1]]
			mu.Unlock()
			if !ok {
				fmt.Println("couldnt find ")
				break
			}
			conn.Write([]byte(value))
		}
		if message[0] == "SET" {
			mu.Lock()
			store[message[1]] = message[2]
			mu.Unlock()
			saving_message := strings.Join(message, " ")

			db.Write([]byte(saving_message + "\n"))

		}
	}

}
func main() {

	load_from_disk()
	listen, err := net.Listen("tcp", ":5379")
	if err != nil {
		fmt.Println("didnt find any connection sorry ")
		return
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("couldnt connect sorry ")
			return
		}
		go makeconnection(conn)
	}
}

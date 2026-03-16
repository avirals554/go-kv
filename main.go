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
func GET(store map[string]string, key string, conn net.Conn) string {
	value, ok := store[key]
	if !ok {
		fmt.Println("there was an error in the get function ")
		conn.Write([]byte("nil"))

	}

	return value

}
func SET(store map[string]string, key string, value string) {
	store[key] = value

}
func makeconnection(conn net.Conn, conn_backup net.Conn) {
	fmt.Println("we are getting a connection from  ", conn.RemoteAddr())
	conn.Write([]byte("tcp connection has started ......." + "\n"))

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("something went went wrong conn.read()")
			break
		}
		message := strings.Split(strings.TrimSpace(string(buf[:n])), " ")
		if message[0] == "GET" {
			mu.Lock()
			value := GET(store, message[1], conn)
			mu.Unlock()
			conn.Write([]byte(value + "\n"))
		}
		if message[0] == "SET" {

			mu.Lock()
			SET(store, message[1], message[2])
			conn.Write([]byte("\n"))
			mu.Unlock()
			saving_message := strings.Join(message, " ") // this is for combining the sliced strings

			db.Write([]byte(saving_message + "\n"))
			if conn_backup != nil {
				conn_backup.Write([]byte(saving_message + "\n"))
			}
		}
	}

}
func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage: go-kv <port> <Leader|Follower>")
		return
	}
	port := os.Args[1]

	load_from_disk()
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("didnt find any connection sorry ")
		return
	}
	var conn_backup net.Conn
	if os.Args[2] == "Leader" {
		conn_backup, err = net.Dial("tcp", "localhost:5380")
		if err != nil {
			fmt.Println("there was a problem connecting to the backup ")
			conn_backup = nil
		}
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("couldnt connect sorry ")
			return
		}
		go makeconnection(conn, conn_backup)
	}
}

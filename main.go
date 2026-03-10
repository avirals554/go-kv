package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	store := make(map[string]string)
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
		fmt.Println("we are getting a connection from  ", conn.RemoteAddr())
		conn.Write([]byte("tcp connection has started ......."))
		for {
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("something went went wrong line 26 onwards ")
				break
			}
			message := strings.Split(string(buf[:n]), " ")
			if message[0] == "GET" {
				value, ok := store[message[1]]
				if !ok {
					fmt.Println("couldnt find ")
					break
				}
				conn.Write([]byte(value))
			}
			if message[0] == "SET" {
				store[message[1]] = message[2]

			}
		}

	}

}

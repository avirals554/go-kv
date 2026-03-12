package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	count := 0

	for {
		conn, err := net.Dial("tcp", "localhost:5379")
		if err != nil {
			count++
		}
		if err == nil {
			conn.Close()
			count = 0
		}
		if count == 3 {
			fmt.Println("the leader has failed ....")
		}
		time.Sleep(3 * time.Second)

	}
}

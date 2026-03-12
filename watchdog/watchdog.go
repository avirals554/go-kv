package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

var message_to_proxy, err = os.OpenFile("active_node.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)

func main() {
	count := 0
	message := "5379"
	message_to_proxy.Truncate(0)
	message_to_proxy.Seek(0, 0)
	message_to_proxy.Write([]byte(message + "\n"))
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
			if message != "5380" {
				message = "5380"
				message_to_proxy.Truncate(0)
				message_to_proxy.Seek(0, 0)
				message_to_proxy.Write([]byte(message + "\n"))
			}
		}
		time.Sleep(3 * time.Second)

	}
}

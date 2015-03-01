// Mud-Client project main.go
package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

type ClientMessage struct {
	Command int
	Value   string
}

type ServerMessage struct {
	Value string
}

func main() {

	service := "127.0.0.1:1200"

	conn, err := net.Dial("tcp", service)
	checkError(err)

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	message := ClientMessage{Command: 1, Value: "test message"}
	fmt.Println("Sending message")

	encoder.Encode(message)
	fmt.Println("message sent")

	var serversResponse ServerMessage
	fmt.Println("waiting for response")
	decoder.Decode(&serversResponse)
	fmt.Println("message received")
	fmt.Println(serversResponse.Value)

	conn.Close()
	os.Exit(0)
}

func printFormattedServerMessage(message string) {

}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

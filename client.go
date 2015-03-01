// Mud-Client project main.go
package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

type ClientMessage struct {
	command int
	value   string
}

type ServerMessage struct {
	value string
}

func main() {

	service := "127.0.0.1:1200"

	conn, err := net.Dial("tcp", service)
	checkError(err)

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	message := ClientMessage{1, "test message"}
	encoder.Encode(message)

	var serversResponse ServerMessage
	decoder.Decode(&serversResponse)
	fmt.Println(serversResponse.value)

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

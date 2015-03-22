// Mud-Client project main.go
package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"sync"
)

type ServerMessage struct {
	Value string
}

var net_lock sync.Mutex

func main() {

	logInTest()
	os.Exit(0)
}

func logInTest() {
	service := "127.0.0.1:1200"

	conn, err := net.Dial("tcp", service)
	checkError(err)

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	message := ClientMessage{Command: 0, Value: "testChar password"}
	fmt.Println("Sending message")

	encoder.Encode(message)
	fmt.Println("message sent")

	var serversResponse ServerMessage
	fmt.Println("waiting for response")
	decoder.Decode(&serversResponse)
	fmt.Println("message received")
	fmt.Println(serversResponse.Value)

	conn.Close()
}

func gobTest() {
	service := "127.0.0.1:1200"

	conn, err := net.Dial("tcp", service)
	checkError(err)

	/*
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
		os.Exit(0)*/

	go receiveMessage(conn)

	encoder := gob.NewEncoder(conn)
	reader := bufio.NewReader(os.Stdin)
	for {

		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		message := ClientMessage{Command: 1, Value: text}
		//net_lock.Lock()
		encoder.Encode(message)
		//net_lock.Unlock()

	}

}

func receiveMessage(conn net.Conn) {

	var serversResponse ServerMessage
	decoder := gob.NewDecoder(conn)
	for {
		//net_lock.Lock()
		err := decoder.Decode(&serversResponse)
		//net_lock.Unlock()
		checkError(err)
		if err == nil {
			fmt.Println("message received")
			fmt.Println(serversResponse.Value)
		}
	}
	conn.Close()
}

func printFormattedServerMessage(message string) {

}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

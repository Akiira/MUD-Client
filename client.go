// Mud-Client project main.go
package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"net"
	"os"
	"sync"
)

type FormattedString struct {
	Color ct.Color
	Value string
}
type ServerMessage struct {
	Value []FormattedString
}

var net_lock sync.Mutex

func main() {

	Test()
	os.Exit(0)
}

func Test() {
	service := "127.0.0.1:1200"

	conn, err := net.Dial("tcp", service)
	checkError(err)

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	message := ClientMessage{Command: "initialMessage", Value: "Ragnar"}
	fmt.Println("Sending message")

	encoder.Encode(message)
	fmt.Println("message sent")

	var serversResponse ServerMessage
	fmt.Println("waiting for response")
	decoder.Decode(&serversResponse)
	fmt.Println("message received: ")
	//fmt.Println(serversResponse.Value)

	printFormatedOutput(serversResponse.Value)

	message = ClientMessage{Command: "say", Value: "test say"}
	encoder.Encode(message)
	fmt.Println("message 2 sent")
	for {
		fmt.Println("waiting for response")
		decoder.Decode(&serversResponse)
		fmt.Println("Message received: ", serversResponse)
		printFormatedOutput(serversResponse.Value)
	}

	conn.Close()
}

func printFormatedOutput(output []FormattedString) {
	for _, element := range output {
		ct.ChangeColor(element.Color, false, ct.Black, false)
		fmt.Println(element.Value)
	}
	ct.ResetColor()
}

func logInTest() {
	service := "127.0.0.1:1200"

	conn, err := net.Dial("tcp", service)
	checkError(err)

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	message := ClientMessage{Command: "look", Value: "testChar password"}
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
		message := ClientMessage{Command: "test", Value: text}
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
func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

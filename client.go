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

	Test2()
	os.Exit(0)
}

func Test2() {
	service := "127.0.0.1:1200"

	conn, err := net.Dial("tcp", service)
	checkError(err)
	defer conn.Close()

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	message := ClientMessage{Command: "initialMessage", Value: "Ragnar"}
	encoder.Encode(message)
	go nonBlockingRead(decoder)
	getInputFromUser(encoder)
}

func getInputFromUser(encoder *gob.Encoder) {
	in := bufio.NewReader(os.Stdin)
	for {
		var msg ClientMessage
		var input string
		read, err := fmt.Scan(&input)
		checkError(err)
		_ = read

		if input == "exit" {
			msg.setToExitMessage()
			break
		} else if input == "attack" {
			var target string
			read, err = fmt.Scan(&target)
			msg.setToAttackMessage(target)
		} else if input == "look" {
			var target string
			read, err = fmt.Scan(&target)
			msg.setToLookMessage(target)
		} else if input == "get" {
			var item string
			read, err = fmt.Scan(&item)
			msg.setToGetMessage(item)
		} else if input == "say" {
			line, err := in.ReadString('\n')
			checkError(err)
			msg.setToSayMessage(line)
		} else { //assume movement
			msg.setToMovementMessage(input)
		}
		fmt.Println("Sending: ", msg)
		encoder.Encode(msg)

		if msg.Command == "exit" {
			break
		}
	}

	os.Exit(0)
}

func nonBlockingRead(decoder *gob.Decoder) {
	for {
		var serversResponse ServerMessage
		decoder.Decode(&serversResponse)
		printFormatedOutput(serversResponse.Value)
	}
}

func printFormatedOutput(output []FormattedString) {
	for _, element := range output {
		ct.ChangeColor(element.Color, false, ct.Black, false)
		fmt.Println(element.Value)
	}
	ct.ResetColor()
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

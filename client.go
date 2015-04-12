// Mud-Client project main.go
package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"net"
	"os"
	"strings"
	"sync"
)

var net_lock sync.Mutex
var conn net.Conn
var encoder *gob.Encoder
var decoder *gob.Decoder
var breakSignal bool

func main() {

	runClient()

	os.Exit(0)
}

func runClient() {
	breakSignal = false
	connectToServer("127.0.0.1:1200") //TODO remove hard coding

	go nonBlockingRead()
	getInputFromUser()
}

func getInputFromUser() {
	in := bufio.NewReader(os.Stdin)
	for {
		var msg ClientMessage
		var input, target string
		read, err := fmt.Scan(&input)
		checkError(err)
		_ = read

		if input == "exit" {
			msg.setToExitMessage()
			break
		} else if input == "attack" {
			read, err = fmt.Scan(&target)
			msg.setToAttackMessage(target)
		} else if input == "look" {
			read, err = fmt.Scan(&target)
			msg.setToLookMessage(target)
		} else if input == "get" {
			var target string
			read, err = fmt.Scan(&target)
			msg.setToGetMessage(target)
		} else if input == "say" {
			line, _ := in.ReadString('\n')
			line = strings.TrimRight(line, "\n")
			msg.setToSayMessage(line)
		} else if input == "stats" {
			msg.setCommand("stats")
		} else if input == "inv" {
			msg.setCommand("inv")
		} else if input == "bid" {
			read, err = fmt.Scan(&target)
			msg.setMsgWithTimestamp("bid", target)
		} else if input == "auction" {
			line, _ := in.ReadString('\n')
			line = strings.TrimRight(line, "\n")
			msg.setCommandAndValue("auction", line)
		} else { //assume movement
			msg.setToMovementMessage(input) //TODO add error handling code here
		}
		fmt.Println("Sending: ", msg)

		net_lock.Lock()
		err = encoder.Encode(msg)
		checkError(err)
		net_lock.Unlock()

		if msg.Command == "exit" {
			break
		}
	}
	breakSignal = true
}

func nonBlockingRead() {
	for {
		var serversResponse ServerMessage
		err := decoder.Decode(&serversResponse)
		checkError(err)

		if serversResponse.MsgType == REDIRECT {
			net_lock.Lock()
			err := conn.Close()
			checkError(err)
			connectToServer(serversResponse.getMessage())
			net_lock.Unlock()
		} else if serversResponse.MsgType == PING {
			net_lock.Lock()
			tmp := newClientMessage("ping", "")
			err := encoder.Encode(&tmp)
			checkError(err)
			net_lock.Unlock()
		} else {
			printFormatedOutput(serversResponse.Value)
		}

		if breakSignal {
			break
		}
	}

	os.Exit(0)
}

func connectToServer(address string) {
	var err error
	fmt.Println("Address:", address)
	conn, err = net.Dial("tcp", address)
	checkError(err)

	encoder = gob.NewEncoder(conn)
	decoder = gob.NewDecoder(conn)

	message := ClientMessage{Command: "initialMessage", Value: "Tiefling password"}
	err = encoder.Encode(message)
	checkError(err)
}

func printFormatedOutput(output []FormattedString) {
	for _, element := range output {
		ct.ChangeColor(element.Color, false, ct.Black, false)
		fmt.Print(element.Value)
	}
	ct.ResetColor()
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

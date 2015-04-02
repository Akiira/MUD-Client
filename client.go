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

var net_lock sync.Mutex
var conn net.Conn
var encoder *gob.Encoder
var decoder *gob.Decoder

func main() {

	//logInTest()

	runClient()

	os.Exit(0)
}

func runClient() {
	connectToServer("127.0.0.1:1200")

	message := ClientMessage{Command: "initialMessage", Value: "Ragnar"}
	encoder.Encode(message)

	go nonBlockingRead()
	getInputFromUser()
}

func getInputFromUser() {
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

		net_lock.Lock()
		encoder.Encode(msg)
		net_lock.Unlock()

		if msg.Command == "exit" {
			break
		}
	}

	os.Exit(0)
}

func nonBlockingRead() {
	for {
		var serversResponse ServerMessage
		decoder.Decode(&serversResponse)

		if serversResponse.MsgType == REDIRECT {
			net_lock.Lock()
			conn.Close()
			connectToServer(serversResponse.getMessage())
			net_lock.Unlock()
		} else {
			printFormatedOutput(serversResponse.Value)
		}
	}
}

func connectToServer(address string) {
	conn, err := net.Dial("tcp", address)
	checkError(err)

	encoder = gob.NewEncoder(conn)
	decoder = gob.NewDecoder(conn)
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

//func logInTest() {
//	service := "127.0.0.1:7200"

//	for {

//		conn, err := net.Dial("tcp", service)
//		checkError(err)
//		fmt.Println("established new connection")

//		encoder := gob.NewEncoder(conn)
//		decoder := gob.NewDecoder(conn)

//		message := ClientMessage{Command: CommandLogin, Value: "Haplo password"}
//		fmt.Println("Sending message")

//		encoder.Encode(message)
//		fmt.Println("message sent")

//		var serversResponse ServerMessage
//		fmt.Println("waiting for response")
//		for {

//			err := decoder.Decode(&serversResponse)
//			if err != nil {
//				fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
//				os.Exit(1)
//			} else {
//				if !serversResponse.isError() {

//					if serversResponse.getServerSystemMessageType() == CommandRedirectServer {
//						service = serversResponse.getServerSystemMessageDetail()
//						break
//					} else if serversResponse.isError() {

//						fmt.Println(serversResponse.getServerSystemMessageType())
//						fmt.Println(serversResponse.getServerSystemMessageDetail())
//						os.Exit(1)
//					}
//				}
//			}

//			conn.Close()
//		}
//	}
//}

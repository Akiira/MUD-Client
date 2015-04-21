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
var username string
var pass string
var serverAddr string
var cacheCharInfo []FormattedString

//var server

func main() {

	readConfigFile()

	runClient()

	os.Exit(0)
}

func readConfigFile() {
	readPassFile, err := os.Open("login.txt")
	checkError(err)

	reader := bufio.NewReader(readPassFile)
	line, _, err := reader.ReadLine()
	username = string(line)
	line, _, err = reader.ReadLine()
	pass = string(line)
	line, _, err = reader.ReadLine()
	serverAddr = string(line)
}

func runClient() {
	breakSignal = false
	connectToServer(serverAddr) //TODO remove hard coding
	go startPingServer()
	go nonBlockingRead()
	NewGetInputFromUser()
}

func startPingServer() {

	listerner := setUpServerWithAddress(":1600")

	for {
		conn, err := listerner.Accept()
		checkError(err)

		pingEnc := gob.NewEncoder(conn)
		pingDec := gob.NewDecoder(conn)

		for {

			var msg ServerMessage

			err := pingDec.Decode(&msg)
			checkError(err)
			fmt.Println("\tPing received.")
			if msg.getMessage() == "done" {
				fmt.Println("Done pinging.")
				break
			}

			err = pingEnc.Encode(&ClientMessage{})
			checkError(err)
		}

		conn.Close()

		if breakSignal {
			break
		}
	}
}

func setUpServerWithAddress(addr string) *net.TCPListener {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	return listener
}

func nonBlockingRead() {
	for {
		var serversResponse ServerMessage
		err := decoder.Decode(&serversResponse)
		//fmt.Println("\tRead message from server.", serversResponse)
		checkError(err)

		if serversResponse.MsgType == REDIRECT {
			net_lock.Lock()
			err := conn.Close()
			checkError(err)
			connectToServer(serversResponse.getMessage())
			net_lock.Unlock()
		} else {
			printFormatedOutput(serversResponse.Value)
			printFormatedOutput(serversResponse.getFormattedCharInfo())
			cacheCharInfo = serversResponse.getFormattedCharInfo()
		}

		if breakSignal {
			break
		}
	}

	os.Exit(0)
}

func NewGetInputFromUser() {
	in := bufio.NewReader(os.Stdin)
	for {
		var msg ClientMessage
		var input string

		_, err := fmt.Scan(&input)
		checkError(err)
		input = strings.TrimSpace(input)

		if isValidDirection(input) {
			msg = newClientMessage("move", input)
		} else if isCombatCommand(input) {
			line, _ := in.ReadString('\n')
			line = strings.TrimSpace(line)
			msg = newClientMessage2(isCombatCommand(input), input, line)
		} else if isNonCombatCommand(input) {
			line, _ := in.ReadString('\n')
			line = strings.TrimSpace(line)
			msg = newClientMessage(input, line)
		} else {
			fmt.Println("\nThat does not appear to be a valid command.\n")
			printFormatedOutput(cacheCharInfo)
			continue
		}

		net_lock.Lock()
		fmt.Println("Sending: ", msg)
		err = encoder.Encode(msg)
		checkError(err)
		net_lock.Unlock()

		if input == "exit" {
			break
		}
	}
	breakSignal = true
}

//TODO consider putting the commands in a file, and then loading them into a hashmap instead
func isLegalCommand(cmd string) bool {
	return isCombatCommand(cmd) || isNonCombatCommand(cmd)
}

func isNonCombatCommand(cmd string) bool {
	switch {
	case cmd == "auction" || cmd == "bid":
		return true
	case cmd == "wield" || cmd == "unwield":
		return true
	case cmd == "equip" || cmd == "unequip":
		return true
	case cmd == "inv":
		return true
	case cmd == "save" || cmd == "exit":
		return true
	case cmd == "stats":
		return true
	case cmd == "look":
		return true
	case cmd == "get" || cmd == "put" || cmd == "drop":
		return true
	case cmd == "move":
		return true
	case cmd == "say" || cmd == "yell":
		return true
	case cmd == "trade":
		return true
	case cmd == "select":
		return true
	case cmd == "reject":
		return true
	case cmd == "accept":
		return true
	case cmd == "help":
		return true

	}
	return isValidDirection(cmd)
}

func isCombatCommand(cmd string) bool {
	switch {
	case cmd == "attack":
		return true
	case cmd == "cast":
		return true
	}

	return false
}

func isValidDirection(direction string) bool {

	switch strings.ToLower(direction) {
	case "n", "n\r\n", "n\n", "north":
		return true
	case "s", "s\r\n", "s\n", "south":
		return true
	case "e", "e\r\n", "e\n", "east":
		return true
	case "w", "w\r\n", "w\n", "west":
		return true
	case "nw", "nw\r\n", "nw\n", "northwest":
		return true
	case "ne", "ne\r\n", "ne\n", "northeast":
		return true
	case "sw", "sw\r\n", "sw\n", "southwest":
		return true
	case "se", "se\r\n", "se\n", "southeast":
		return true
	case "u", "u\r\n", "u\n", "up":
		return true
	case "d", "d\r\n", "d\n", "down":
		return true
	}

	return false
}

func connectToServer(address string) {
	var err error
	fmt.Println("\tAddress:", address)
	conn, err = net.Dial("tcp", address)
	fmt.Println("\tFinished dialing.")
	checkError(err)

	encoder = gob.NewEncoder(conn)
	decoder = gob.NewDecoder(conn)

	message := ClientMessage{Command: "initialMessage", Value: username + " " + pass} //TODO get this from user
	err = encoder.Encode(&message)
	fmt.Println("\tMessage Sent")
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

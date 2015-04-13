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
	go startPingServer()
	go nonBlockingRead()
	//OldGetInputFromUser()
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

func OldGetInputFromUser() {
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
			line = strings.TrimSpace(line)
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
			line = strings.TrimSpace(line)
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
		} else {
			printFormatedOutput(serversResponse.Value)
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

		line, err := in.ReadString('\n')
		checkError(err)
		line = strings.TrimSpace(line)

		if isCombatCommand(input) { //TODO make this if-chain a little smaller
			msg = newClientMessage2(isCombatCommand(input), input, line)
		} else if isValidDirection(input) {
			msg = newClientMessage("move", input)
		} else if isNonCombatCommand(input) {
			msg = newClientMessage(input, line)
		} else {
			fmt.Println("\nThat does not appear to be a valid command.\n")
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

//TODO consider putting this in a file, and then loading them into a hashmap instead
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
	fmt.Println("Address:", address)
	conn, err = net.Dial("tcp", address)
	checkError(err)

	encoder = gob.NewEncoder(conn)
	decoder = gob.NewDecoder(conn)

	message := ClientMessage{Command: "initialMessage", Value: "Tiefling password"} //TODO get this from user
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

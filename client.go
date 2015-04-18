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

var (
	net_lock    sync.Mutex
	conn        net.Conn
	encoder     *gob.Encoder
	decoder     *gob.Decoder
	breakSignal bool
	username    string
	password    string
	loggedIn    bool
)

func main() {

	runClient()

	os.Exit(0)
}

func runClient() {
	breakSignal = false

	choice := Start()

	if choice == "new" {
		CreateNewCharacter()
		LogInAndPlay()
	} else if choice == "login" {
		LogInAndPlay()
	} else {
		fmt.Println("That was not a recognized command, goodbye.")
	}
}

func Start() (choice string) {
	fmt.Println("Hail travaler. Would you like to create a new adventerer or login in with an old one?")
	fmt.Println("Type 'new' or 'login'")
	fmt.Scanln(&choice)

	return choice
}

func CreateNewCharacter() {
	connectToServer("127.0.0.1:1202")
	go ReadFromServer()

	for {
		var input string
		fmt.Scan(&input)
		encoder.Encode(newClientMessage("", input))

		if input == "done" { //TODO make this better
			break
		}
	}

	for { //TODO make this better
		if breakSignal {
			break
		}
	}

	breakSignal = false
}

func LogInAndPlay() {
	connectToServer("127.0.0.1:1200") //TODO remove hard coding

	fmt.Println("Please enter you adventuers name.")
	_, err := fmt.Scan(&username)
	checkError(err)
	fmt.Println("Please enter your password.")
	_, err = fmt.Scan(&password)
	checkError(err)

	message := ClientMessage{Command: "initialMessage", Value: username + " " + password}
	err = encoder.Encode(&message)
	checkError(err)

	loggedIn = true

	go ReadFromServer()
	GetInputFromUser()
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

			if msg.getMessage() == "done" {
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

func ReadFromServer() {
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
			printFormatedOutput(serversResponse.getFormattedCharInfo())
		}

		if breakSignal || serversResponse.MsgType == EXIT {
			breakSignal = true
			break
		}
	}
}

func GetInputFromUser() {
	in := bufio.NewReader(os.Stdin)
	for {
		var msg ClientMessage
		var input string

		_, err := fmt.Scan(&input)
		checkError(err)
		input = strings.TrimSpace(input)

		if isValidDirection(input) {
			msg = newClientMessage("move", input)
		} else if input == "accept" || input == "done" {
			msg = newClientMessage(input, input)
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
	case cmd == "opentrade":
		return true
	case cmd == "add":
		return true
	case cmd == "accept":
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

	fmt.Println("\tAddress: ", address)
	conn, err = net.Dial("tcp", address)
	fmt.Println("\tFinished dialing.")
	checkError(err)

	encoder = gob.NewEncoder(conn)
	decoder = gob.NewDecoder(conn)

	if loggedIn {
		encoder.Encode(newClientMessage("initialMessage", username+" "+password))
	}
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

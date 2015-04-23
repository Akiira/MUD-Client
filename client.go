// Mud-Client project main.go
package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	servers     map[string]string = make(map[string]string)
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
	readServerList()
	runClient()

	os.Exit(0)
}

func readServerList() {
	file, err := os.Open("serverList.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		readData := strings.Fields(scanner.Text())
		fmt.Println(readData)
		servers[readData[0]] = readData[1]
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
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
	connectToServer(servers["newChar"])
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

	connectToServer(servers["central"])

	fmt.Println("Please enter you adventuers name.")
	_, err := fmt.Scan(&username)
	checkError(err)
	fmt.Println("Please enter your password.")
	_, err = fmt.Scan(&password)
	checkError(err)

	err = encoder.Encode(&ClientMessage{Command: "initialMessage", Value: username + " " + password})
	checkError(err)

	loggedIn = true

	go ReadFromServer()
	GetInputFromUser()
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
		} else if serversResponse.MsgType == PING {
			encoder.Encode(newClientMessage("ping", "ping"))
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
		} else if isOneWordCommand(input) {
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

func isOneWordCommand(cmd string) bool {
	switch cmd {
	case "accept", "done":
		return true
	case "equipment", "eq":
		return true
	case "stats":
		return true
	case "look":
		return true
	}
	return false
}

//TODO consider putting the commands in a file, and then loading them into a hashmap instead
func isLegalCommand(cmd string) bool {
	return isCombatCommand(cmd) || isNonCombatCommand(cmd)
}

func isNonCombatCommand(cmd string) bool {
	switch cmd {
	case "auction", "bid":
		return true
	case "wield", "unwield":
		return true
	case "equip", "unequip", "wear", "remove":
		return true
	case "inv", "inventory", "eq", "equipment":
		return true
	case "save", "exit":
		return true
	case "stats":
		return true
	case "look":
		return true
	case "get", "put", "drop":
		return true
	case "move", "flee":
		return true
	case "say", "yell":
		return true
	case "trade", "opentrade", "add":
		return true
	case "accept":
		return true
	}
	return isValidDirection(cmd)
}

func isCombatCommand(cmd string) bool {
	switch cmd {
	case "attack", "a", "kill":
		return true
	case "cast":
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
	var err error //Dont remove me

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

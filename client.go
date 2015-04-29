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
	servers         map[string]string = make(map[string]string)
	net_lock        sync.Mutex
	conn            net.Conn
	encoder         *gob.Encoder
	decoder         *gob.Decoder
	breakSignal     bool
	username        string
	password        string
	loggedIn        bool
	cachePlayerInfo []FormattedString
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
		//fmt.Println(readData)
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

func connectToServer(address string) {
	var err error //Dont remove me

	conn, err = net.Dial("tcp", address)
	checkError(err)

	encoder = gob.NewEncoder(conn)
	decoder = gob.NewDecoder(conn)

	if loggedIn {
		encoder.Encode(newClientMessage("initialMessage", username+" "+password))
	}
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
			cachePlayerInfo = serversResponse.getFormattedCharInfo()
			printFormatedOutput(cachePlayerInfo)
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

		line, err := in.ReadString('\n')
		checkError(err)
		line = strings.TrimSpace(strings.ToLower(line))

		cmd := strings.Split(line, " ")[0]

		if isValidDirection(cmd) {
			msg = newClientMessage("move", cmd)
		} else {
			msg = newClientMessage2(isCombatCommand(cmd), cmd, line)
		}

		net_lock.Lock()
		//fmt.Println("Sending: ", msg)
		err = encoder.Encode(msg)
		checkError(err)
		net_lock.Unlock()

		if cmd == "exit" {
			break
		}
	}
	breakSignal = true
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

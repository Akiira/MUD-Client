// ClientMessage
package main

import (
	"strings"
	"time"
)

//this is suppose to be an event
type ClientMessage struct {
	CombatAction bool
	Command      string
	Value        string
}

func newClientMessage(cmd string, val string) ClientMessage {
	return ClientMessage{CombatAction: false, Command: cmd, Value: val}
}
func newClientMessage2(isCombat bool, cmd string, line string) ClientMessage {
	val := strings.TrimSpace(strings.TrimPrefix(line, cmd))
	
	msg := ClientMessage{CombatAction: isCombat, Command: cmd, Value: val}
	
	if msg.IsTradeCommand() && len(val) == 0 {
		msg.setValue(cmd)
	}
	
	return msg
}

func (msg *ClientMessage) IsTradeCommand() bool {
	switch msg.Command {
	case "accept", "done", "add", "reject":
		return true
	}

	return false
}

func (msg *ClientMessage) setValue(val string) {
	msg.Value = val
}

func (msg *ClientMessage) setCommand(cmd string) {
	msg.Command = cmd
}

func (msg *ClientMessage) setCommandAndValue(cmd string, val string) {
	msg.Command = cmd
	msg.Value = val
}

func (msg *ClientMessage) setCommandWithTimestamp(cmd string) {
	msg.Command = cmd + ";" + time.Now().String()
}

func (msg *ClientMessage) setMsgWithTimestamp(cmd string, value string) {
	msg.Command = cmd + ";" + time.Now().String()
	msg.Value = value
}

func (msg *ClientMessage) getTimeStamp() string {

	peices := strings.Split(msg.Command, ";")
	if len(peices) == 2 {
		return peices[1]
	} else {
		return ""
	}
}


func (message *ClientMessage) getPassword() string {
	split := strings.Split(message.Value, " ")
	return split[1]
}

func (message *ClientMessage) getUsername() string {
	split := strings.Split(message.Value, " ")
	return split[0]
}

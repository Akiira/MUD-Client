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
	return ClientMessage{CombatAction: isCombat, Command: cmd, Value: val}
}

func (msg *ClientMessage) setCommand(cmd string) {
	msg.CombatAction = false
	msg.Command = cmd
	msg.Value = ""
}

func (msg *ClientMessage) setCommandAndValue(cmd string, val string) {
	msg.CombatAction = false
	msg.Command = cmd
	msg.Value = val
}

func (msg *ClientMessage) setCommandWithTimestamp(cmd string) {
	msg.CombatAction = false
	msg.Command = cmd + ";" + time.Now().String()
	msg.Value = ""
}

func (msg *ClientMessage) setMsgWithTimestamp(cmd string, value string) {
	msg.CombatAction = false
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

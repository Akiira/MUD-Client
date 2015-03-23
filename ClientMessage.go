// ClientMessage
package main

import (
	"strings"
)

type ClientMessage struct {
	combatAction bool
	Command      string
	Value        string
}

func ClientMessageConstructor(cmd string, val string) ClientMessage {
	return ClientMessage{combatAction: false, Command: cmd, Value: val}
}

func (msg *ClientMessage) setToMovementMessage(direction string) {
	msg.combatAction = false
	msg.Command = "move"
	msg.Value = direction
}

func (msg *ClientMessage) setToGetMessage(item string) {
	msg.combatAction = false
	msg.Command = "get"
	msg.Value = item
}

func (msg *ClientMessage) setToLookMessage(target string) {
	msg.combatAction = false
	msg.Command = "look"
	msg.Value = target
}

func (msg *ClientMessage) setToAttackMessage(target string) {
	msg.combatAction = true
	msg.Command = "attack"
	msg.Value = target
}

func (msg *ClientMessage) setToExitMessage() {
	msg.combatAction = false
	msg.Command = "exit"
	msg.Value = ""
}

func (msg *ClientMessage) setAll(combatAction bool, cmd string, val string) {
	msg.combatAction = combatAction
	msg.Command = cmd
	msg.Value = val
}

func (msg *ClientMessage) setAllNonCombat(cmd string, val string) {
	msg.combatAction = false
	msg.Command = cmd
	msg.Value = val
}

func (message *ClientMessage) getPassword() string {
	split := strings.Split(message.Value, " ")
	return split[1]
}

func (message *ClientMessage) getUsername() string {
	split := strings.Split(message.Value, " ")
	return split[0]
}

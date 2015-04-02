// ClientMessage
package main

import (
	"strings"
)

////command for error
//const ServerErrorMessage = "error"
//const ErrorUnexpectedCommand = 201
//const ErrorWorldIsNotFound = 202
//const ErrorAuthorizationFail = 203

////command for system
//const CommandLogin = "Login"
//const CommandLogout = 102
//const CommandRedirectServer = "RedirectServer"
//const CommandEnterWorld = 104
//const CommandQueryCharacter = "QueryCharacter"
//const CommandSaveCharacter = "SaveCharacter"

////command for create user
//const CommandRegister = 111

////command in a room
//const CommandAttack = 11
//const CommandItem = 12
//const CommandLeave = 13 // leave occur the same time with enter the room??

////command between room?
//const CommandJoinWorld = 21 // will change the room occur the same time with leave?
//// probably use after authenticate with login server and move to the first world as well

//this is suppose to be an event
type ClientMessage struct {
	CombatAction bool
	Command      string
	Value        string
}

func ClientMessageConstructor(cmd string, val string) ClientMessage {
	return ClientMessage{CombatAction: false, Command: cmd, Value: val}
}

func (msg *ClientMessage) setToMovementMessage(direction string) {
	msg.CombatAction = false
	msg.Command = "move"
	msg.Value = direction
}

func (msg *ClientMessage) setToSayMessage(thingToSay string) {
	msg.CombatAction = false
	msg.Command = "say"
	msg.Value = thingToSay
}

func (msg *ClientMessage) setToGetMessage(item string) {
	msg.CombatAction = false
	msg.Command = "get"
	msg.Value = item
}

func (msg *ClientMessage) setToLookMessage(target string) {
	msg.CombatAction = false
	msg.Command = "look"
	msg.Value = target
}

func (msg *ClientMessage) setToAttackMessage(target string) {
	msg.CombatAction = true
	msg.Command = "attack"
	msg.Value = target
}

func (msg *ClientMessage) setToExitMessage() {
	msg.CombatAction = false
	msg.Command = "exit"
	msg.Value = ""
}

func (msg *ClientMessage) setAll(combatAction bool, cmd string, val string) {
	msg.CombatAction = combatAction
	msg.Command = cmd
	msg.Value = val
}

func (msg *ClientMessage) setAllNonCombat(cmd string, val string) {
	msg.CombatAction = false
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

package main

import (
	"./server"
	"./utils"
)

func main() {
	// initializes state info
	var state = initialize()

	var connection = server.NewConnection()
	connection.Register(&state)
}

func initialize() utils.State {
	state := utils.State{}

	state.Run = -1

	// NOTE: for sake of simplicity assuming user would have only 1 message to send
	// although this project can handle multi message clients also
	state.MyMsgCount = 1
	state.MyMessages = make([]string, state.MyMsgCount)
	state.MyMessagesHash = make([]uint64, state.MyMsgCount)

	// generate random 160 bit message
	// NOTE: assuming that all peers would have only one message
	state.MyMessages[0] = utils.GenerateMessage()

	return state
}

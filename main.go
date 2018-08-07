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
	state.MyMsgCount = 1
	state.MyMessages = make([]string, state.MyMsgCount)
	state.MyMessagesHash = make([]uint64, state.MyMsgCount)

	// generate random 160 bit message
	// NOTE: assuming that all peers would have only one message
	state.MyMessages[0] = utils.GenerateMessage()

	return state
}

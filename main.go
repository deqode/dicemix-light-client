package main

import (
	"math/rand"
	"time"

	"./server"
	"./utils"
)

func main() {
	// TODO: generate LTSK, LTPK
	// Broadcast LTPK before KEPK
	// sign all messages from LTSK

	// initializes state info
	var state = initialize()

	var connection = server.NewConnection()
	connection.Register(&state)
}

func initialize() utils.State {
	state := utils.State{}

	state.Run = -1

	// NOTE: for sake of simplicity assuming user would generate random n messages
	// 0 < n < 4
	state.MyMsgCount = count()
	state.MyMessages = make([]string, state.MyMsgCount)
	state.MyMessagesHash = make([]uint64, state.MyMsgCount)

	// generate random 160 bit message
	for i := 0; i < int(state.MyMsgCount); i++ {
		state.MyMessages[i] = utils.GenerateMessage()
	}

	return state
}

// return randomly generated n
// 0 < n < 4
func count() uint32 {
	rand.Seed(time.Now().UnixNano())
	return uint32(rand.Intn(3-1) + 1)
}

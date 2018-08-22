package main

import (
	"math/rand"
	"time"

	"./eddsa"
	"./server"
	"./utils"
	log "github.com/sirupsen/logrus"
)

func main() {
	// setup logger
	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)

	// initializes state info
	var state = initialize()

	log.Info("Attempt to connect to DiceMix Server")
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

	edDSA := eddsa.NewCurveED25519()
	state.Ltpk, state.Ltsk, _ = edDSA.GenerateKeyPair()

	return state
}

// return randomly generated n
// 0 < n < 4
func count() uint32 {
	rand.Seed(time.Now().UnixNano())
	return uint32(rand.Intn(3-1) + 1)
}

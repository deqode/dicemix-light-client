package main

import (
	"math/rand"
	"time"

	"github.com/techracers-blockchain/dicemix-light-client/ecdsa"

	"github.com/techracers-blockchain/dicemix-light-client/server"
	"github.com/techracers-blockchain/dicemix-light-client/utils"

	log "github.com/sirupsen/logrus"
)

// Entry point
func main() {
	// setup logger
	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)

	// initializes state info
	var state = initialize()

	log.Info("Attempt to connect to DiceMix Server")

	// creating a new websocket connection with server
	var connection = server.NewConnection()
	connection.Register(&state)
}

func initialize() utils.State {
	state := utils.State{}

	// NOTE: for sake of simplicity assuming user would generate random n messages
	// 0 < n < 4
	state.MyMsgCount = count()
	state.MyMessages = make([]string, state.MyMsgCount)
	state.MyMessagesHash = make([]uint64, state.MyMsgCount)

	// generate my LTSK, LTPK
	ecdsa := ecdsa.NewCurveECDSA()
	state.Session.Ltpk, state.Session.Ltsk, _ = ecdsa.GenerateKeyPair()

	return state
}

// return randomly generated n
// 0 < n < 4
// NOTE: in actual implementation this should return count of your mesages
func count() uint32 {
	rand.Seed(time.Now().UnixNano())
	return uint32(rand.Intn(3-1) + 1)
}

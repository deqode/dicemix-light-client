package server

import (
	"time"

	"../ecdh"
	"../messages"
	"../utils"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// identifies response message from server
// and passes response to appropriate handle for further operations
func handleMessage(conn *websocket.Conn, message []byte, code uint32, state *utils.State) {
	log.WithFields(log.Fields{
		"code": code,
	}).Info("RECV:")

	switch code {
	case messages.S_JOIN_RESPONSE:
		// Response against request to join dicemix transaction
		response := &messages.RegisterResponse{}
		err := proto.Unmarshal(message, response)
		checkError(err)
		handleJoinResponse(response, state)
	case messages.S_START_DICEMIX:
		// Response to start DiceMix Run
		response := &messages.DiceMixResponse{}
		err := proto.Unmarshal(message, response)
		checkError(err)
		handleStartDicemix(conn, response, state)
	case messages.S_KEY_EXCHANGE:
		// Response against request for KeyExchange
		response := &messages.DiceMixResponse{}
		err := proto.Unmarshal(message, response)
		checkError(err)
		handleKeyExchangeResponse(conn, response, state)
	case messages.S_EXP_DC_VECTOR:
		// contains roots of DC-Combined
		response := &messages.DCExpResponse{}
		err := proto.Unmarshal(message, response)
		checkError(err)
		handleDCExpResponse(conn, response, state)
	case messages.S_SIMPLE_DC_VECTOR:
		// conatins peers DC-SIMPLE-VECTOR's
		response := &messages.DiceMixResponse{}
		err := proto.Unmarshal(message, response)
		checkError(err)
		handleDCSimpleResponse(conn, response, state)
	case messages.S_TX_SUCCESSFUL:
		// conatins success message for TX
		response := &messages.TXDoneResponse{}
		err := proto.Unmarshal(message, response)
		checkError(err)
		handleTXDoneResponse(conn, response, state)
	case messages.S_KESK_REQUEST:
		response := &messages.InitiaiteKESK{}
		err := proto.Unmarshal(message, response)
		checkError(err)
		handleKESKRequest(conn, response, state)
	}
}

// Response against request to join dicemix transaction
func handleJoinResponse(response *messages.RegisterResponse, state *utils.State) {
	if response.Err != "" {
		log.Fatal("Error- ", response.Err)
	}
	// stores MyId provided by user
	state.MyID = response.Id

	log.Info(response.Message)
	log.Info("My Id - ", state.MyID)
}

// Response to start DiceMix Run
func handleStartDicemix(conn *websocket.Conn, response *messages.DiceMixResponse, state *utils.State) {
	if response.Err != "" {
		log.Fatal("Error - ", response.Err)
	}

	log.Info("DiceMix protocol has been initiated")

	// initialize variables
	state.Peers = make([]utils.Peers, len(response.Peers)-1)
	set := make(map[int32]struct{}, len(response.Peers)-1)
	i := 0

	// store peers ID's
	// check for duplicate peer id's
	for _, peer := range response.Peers {
		if _, ok := set[peer.Id]; ok {
			log.Fatal("Duplicate peer ID's: ", peer.Id)
		}

		set[peer.Id] = struct{}{}

		if peer.Id != state.MyID {
			state.Peers[i].ID = peer.Id
			i++
		}
	}

	log.Info("Number of peers - ", len(state.Peers))

	// generates NIKE KeyPair for current run
	// mode = 0 to generate (my_kesk, my_kepk)
	iNike.GenerateKeys(state, 0)

	log.Info("MY KESK - ", state.Kesk)
	log.Info("MY KEPK - ", state.Kepk)

	// KeyExchange
	// broadcast our NIKE PublicKey with our peers
	ecdh := ecdh.NewCurve25519ECDH()
	keyExchangeRequest, err := proto.Marshal(&messages.KeyExchangeRequest{
		Code:      messages.C_KEY_EXCHANGE,
		Id:        state.MyID,
		PublicKey: ecdh.Marshal(state.Kepk),
		NumMsgs:   state.MyMsgCount,
		Timestamp: timestamp(),
	})

	// broadcast our PublicKey
	broadcast(conn, keyExchangeRequest, err, messages.C_KEY_EXCHANGE)
}

// Response against request for KeyExchange
func handleKeyExchangeResponse(conn *websocket.Conn, response *messages.DiceMixResponse, state *utils.State) {
	if response.Err != "" {
		log.Fatal("Error - ", response.Err)
	}

	// generate random 160 bit message
	for i := 0; i < int(state.MyMsgCount); i++ {
		state.MyMessages[i] = utils.GenerateMessage()
	}

	log.Info("My Message - ", utils.Base58StringToBytes(state.MyMessages[0]))

	// copies peers info returned from server to local state.Peers
	// store peers PublicKey and NumMsgs
	filterPeers(state, response.Peers)

	// derive shared keys with peers
	iNike.DeriveSharedKeys(state)

	// generate DC Exponential Vector
	iDcNet.DeriveMyDCVector(state)

	// DC EXP
	// broadcast our DC-EXP vector with peers
	dcExpRequest, err := proto.Marshal(&messages.DCExpRequest{
		Code:        messages.C_EXP_DC_VECTOR,
		Id:          state.MyID,
		DCExpVector: state.MyDC,
		Timestamp:   timestamp(),
	})

	// broadcast our my_dc[]
	broadcast(conn, dcExpRequest, err, messages.C_EXP_DC_VECTOR)
}

// obtains roots and runs DC_SIMPLE
func handleDCExpResponse(conn *websocket.Conn, response *messages.DCExpResponse, state *utils.State) {
	if response.Err != "" {
		log.Fatal("Error - ", response.Err)
	}

	// store roots (message hashes) calculated by server
	state.AllMsgHashes = response.Roots

	log.Info("RECV: Roots - ", state.AllMsgHashes)

	// run a SIMPLE DC NET
	iDcNet.RunDCSimple(state)

	if state.NextKepk == nil {
		// generates NIKE KeyPair for next run
		// mode = 1 to generate (my_next_kesk, my_next_kepk)
		iNike.GenerateKeys(state, 1)
	}

	// broadcast our DC SIMPLE Vector
	ecdh := ecdh.NewCurve25519ECDH()
	dcSimpleRequest, err := proto.Marshal(&messages.DCSimpleRequest{
		Code:           messages.C_SIMPLE_DC_VECTOR,
		Id:             state.MyID,
		DCSimpleVector: state.DCSimpleVector,
		MyOk:           state.MyOk,
		NextPublicKey:  ecdh.Marshal(state.NextKepk),
		Timestamp:      timestamp(),
	})

	broadcast(conn, dcSimpleRequest, err, messages.C_SIMPLE_DC_VECTOR)
}

// handles other peers DC-SIMPLE-VECTORS
// resolves DC-NET
func handleDCSimpleResponse(conn *websocket.Conn, response *messages.DiceMixResponse, state *utils.State) {
	if response.Err != "" {
		log.Fatal("Error - ", response.Err)
	}

	// copies peers info returned from server to local state.Peers
	// store other peers DC Simple Vectors
	filterPeers(state, response.Peers)

	// finally resolves DC Net Vectors to obtain messages
	// should contain all honest peers messages in absence of malicious peers
	iDcNet.ResolveDCNet(state)

	// Verify that every peer agrees to proceed
	ok := iDcNet.VerifyProceed(state)

	log.Info("Agree to Proceed? = ", ok)

	var confirmation = make([]byte, 0)
	if ok {
		ecdh := ecdh.NewCurve25519ECDH()
		confirmation, _ = ecdh.GenerateConfirmation(state.AllMessages, state.Kesk)
	}

	// broadcast our Confirmation
	confirmationRequest, err := proto.Marshal(&messages.ConfirmationRequest{
		Code:         messages.C_TX_CONFIRMATION,
		Id:           state.MyID,
		Confirmation: confirmation,
		Messages:     state.AllMessages,
		Timestamp:    timestamp(),
	})

	broadcast(conn, confirmationRequest, err, messages.C_TX_CONFIRMATION)
}

// handles success message for TX
func handleTXDoneResponse(conn *websocket.Conn, response *messages.TXDoneResponse, state *utils.State) {
	if response.Err != "" {
		log.Fatal("Error - ", response.Err)
	}

	log.Info("Transaction successful. All peers agreed.")
	conn.Close()
}

// handles request from server to initiate kesk
// broadcasts our kesk for current round
// initializes - (my_kesk, my_kepk) := (my_next_kesk, my_next_kepk)
func handleKESKRequest(conn *websocket.Conn, response *messages.InitiaiteKESK, state *utils.State) {
	if response.Err != "" {
		log.Fatal("Error - ", response.Err)
	}

	log.Info("RECV: ", response.Message)

	// broadcast our kesk
	ecdh := ecdh.NewCurve25519ECDH()
	initiaiteKESK, err := proto.Marshal(&messages.InitiaiteKESKResponse{
		Code:       messages.C_KESK_RESPONSE,
		Id:         state.MyID,
		PrivateKey: ecdh.MarshalSK(state.Kesk),
		Timestamp:  timestamp(),
	})

	// broadcast our kesk
	broadcast(conn, initiaiteKESK, err, messages.C_KESK_RESPONSE)

	// Rotate keys
	state.Kesk = state.NextKesk
	state.Kepk = state.NextKepk

	state.NextKesk = nil
	state.NextKepk = nil
}

// send request to server
func broadcast(conn *websocket.Conn, request []byte, err error, code int) {
	checkError(err)
	err = conn.WriteMessage(websocket.BinaryMessage, request)
	checkError(err)

	log.WithFields(log.Fields{
		"code": code,
	}).Info("SENT: ")
}

// to identify time of occurence of an event
// returns current timestamp
// example - 2018-08-07 12:04:46.456601867 +0000 UTC m=+0.000753626
func timestamp() string {
	return time.Now().String()
}

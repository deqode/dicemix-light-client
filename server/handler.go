package server

import (
	"github.com/manjeet-thadani/dicemix-client/ecdh"
	"github.com/manjeet-thadani/dicemix-client/messages"
	"github.com/manjeet-thadani/dicemix-client/utils"

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
		handleJoinResponse(conn, response, state)
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
func handleJoinResponse(conn *websocket.Conn, response *messages.RegisterResponse, state *utils.State) {
	if response.Header.Err != "" {
		log.Fatal("Error- ", response.Header.Err)
	}
	// stores MyId provided by user
	state.Session.MyID = response.Id

	log.Info("MY Ltsk - ", state.Session.Ltsk)
	log.Info("MY Ltpk - ", state.Session.Ltpk)

	log.Info(response.Header.Message)
	log.Info("My Id - ", state.Session.MyID)

	// create proto to send response against S_JOIN_RESPONSE
	header := requestHeader(messages.C_LTPK_REQUEST, state.Session.SessionID, state.Session.MyID)
	message, _ := proto.Marshal(&messages.LtpkExchangeRequest{
		Header:    header,
		PublicKey: state.Session.Ltpk,
	})

	// cannot sign message via actual ltsk
	// as uptill now server not have our ltpk
	ltpkExchangeRequest, err := proto.Marshal(&messages.SignedRequest{
		RequestData: message,
		Signature:   []byte{},
	})

	// send our Long Term PublicKey
	send(conn, ltpkExchangeRequest, err, messages.C_LTPK_REQUEST)
}

// Response to start DiceMix Run
func handleStartDicemix(conn *websocket.Conn, response *messages.DiceMixResponse, state *utils.State) {
	if response.Header.Err != "" {
		log.Fatal("Error - ", response.Header.Err)
	}

	log.Info("DiceMix protocol has been initiated")

	// initialize variables
	state.Session.SessionID = response.Header.SessionId
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

		if peer.Id != state.Session.MyID {
			state.Peers[i].ID = peer.Id
			i++
		}
	}

	log.Info("Session Id - ", state.Session.SessionID)
	log.Info("Number of peers - ", len(state.Peers))

	// generates NIKE KeyPair for current run
	// mode = 0 to generate (my_kesk, my_kepk)
	iNike.GenerateKeys(state, 0)

	log.Info("MY KESK - ", state.Session.Kesk)
	log.Info("MY KEPK - ", state.Session.Kepk)

	// KeyExchange
	// send our NIKE PublicKey to server
	header := requestHeader(messages.C_KEY_EXCHANGE, state.Session.SessionID, state.Session.MyID)

	ecdh := ecdh.NewCurve25519ECDH()
	message, err := proto.Marshal(&messages.KeyExchangeRequest{
		Header:    header,
		PublicKey: ecdh.Marshal(state.Session.Kepk),
		NumMsgs:   state.MyMsgCount,
	})

	// generate signed message using our ltsk
	keyExchangeRequest, err := generateSignedRequest(state.Session.Ltsk, message)

	// send our PublicKey
	send(conn, keyExchangeRequest, err, messages.C_KEY_EXCHANGE)
}

// Response against request for KeyExchange
func handleKeyExchangeResponse(conn *websocket.Conn, response *messages.DiceMixResponse, state *utils.State) {
	if response.Header.Err != "" {
		log.Fatal("Error - ", response.Header.Err)
	}

	// generate random 160 bit message
	for i := 0; i < int(state.MyMsgCount); i++ {
		state.MyMessages[i] = utils.GenerateMessage()
	}

	log.Info("My Message (1) - ", utils.Base58StringToBytes(state.MyMessages[0]))

	// copies peers info returned from server to local state.Peers
	// store peers PublicKey and NumMsgs
	filterPeers(state, response.Peers)

	// derive shared keys with peers
	iNike.DeriveSharedKeys(state)

	// generate DC Exponential Vector
	iDcNet.DeriveMyDCVector(state)

	// DC EXP
	// send our DC-EXP vector with peers
	header := requestHeader(messages.C_EXP_DC_VECTOR, state.Session.SessionID, state.Session.MyID)
	message, err := proto.Marshal(&messages.DCExpRequest{
		Header:      header,
		DCExpVector: state.MyDC,
	})

	// generate signed message using our ltsk
	dcExpRequest, err := generateSignedRequest(state.Session.Ltsk, message)

	// send our my_dc[]
	send(conn, dcExpRequest, err, messages.C_EXP_DC_VECTOR)
}

// obtains roots and runs DC_SIMPLE
func handleDCExpResponse(conn *websocket.Conn, response *messages.DCExpResponse, state *utils.State) {
	if response.Header.Err != "" {
		log.Fatal("Error - ", response.Header.Err)
	}

	// store roots (message hashes) calculated by server
	state.AllMsgHashes = response.Roots

	log.Info("RECV: Roots - ", state.AllMsgHashes)

	// run a SIMPLE DC NET
	iDcNet.RunDCSimple(state)

	if state.Session.NextKepk == nil {
		// generates NIKE KeyPair for next run
		// mode = 1 to generate (my_next_kesk, my_next_kepk)
		iNike.GenerateKeys(state, 1)
	}

	// send our DC SIMPLE Vector
	header := requestHeader(messages.C_SIMPLE_DC_VECTOR, state.Session.SessionID, state.Session.MyID)

	ecdh := ecdh.NewCurve25519ECDH()
	message, err := proto.Marshal(&messages.DCSimpleRequest{
		Header:         header,
		DCSimpleVector: state.DCSimpleVector,
		MyOk:           state.MyOk,
		NextPublicKey:  ecdh.Marshal(state.Session.NextKepk),
	})

	// generate signed message using our ltsk
	dcSimpleRequest, err := generateSignedRequest(state.Session.Ltsk, message)

	send(conn, dcSimpleRequest, err, messages.C_SIMPLE_DC_VECTOR)
}

// handles other peers DC-SIMPLE-VECTORS
// resolves DC-NET
func handleDCSimpleResponse(conn *websocket.Conn, response *messages.DiceMixResponse, state *utils.State) {
	if response.Header.Err != "" {
		log.Fatal("Error - ", response.Header.Err)
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

	var confirmation = ok

	// send our Confirmation
	header := requestHeader(messages.C_TX_CONFIRMATION, state.Session.SessionID, state.Session.MyID)
	message, err := proto.Marshal(&messages.ConfirmationRequest{
		Header:       header,
		Confirmation: confirmation,
		Messages:     state.AllMessages,
	})

	// generate signed message using our ltsk
	confirmationRequest, err := generateSignedRequest(state.Session.Ltsk, message)

	send(conn, confirmationRequest, err, messages.C_TX_CONFIRMATION)
}

// handles success message for TX
func handleTXDoneResponse(conn *websocket.Conn, response *messages.TXDoneResponse, state *utils.State) {
	if response.Header.Err != "" {
		log.Fatal("Error - ", response.Header.Err)
	}

	// transaction is successfull
	// close the connection
	log.Info("Transaction successful. All peers agreed.")
	conn.Close()
}

// handles request from server to initiate kesk
// sends our kesk for current round
// initializes - (my_kesk, my_kepk) := (my_next_kesk, my_next_kepk)
func handleKESKRequest(conn *websocket.Conn, response *messages.InitiaiteKESK, state *utils.State) {
	if response.Header.Err != "" {
		log.Fatal("Error - ", response.Header.Err)
	}

	// request to send our KESK to initiate blame stage
	log.Info("RECV: ", response.Header.Message)

	// send our kesk
	header := requestHeader(messages.C_KESK_RESPONSE, state.Session.SessionID, state.Session.MyID)

	ecdh := ecdh.NewCurve25519ECDH()
	message, err := proto.Marshal(&messages.InitiaiteKESKResponse{
		Header:     header,
		PrivateKey: ecdh.MarshalSK(state.Session.Kesk),
	})

	// generate signed message using our ltsk
	initiaiteKESK, err := generateSignedRequest(state.Session.Ltsk, message)

	// send our kesk
	send(conn, initiaiteKESK, err, messages.C_KESK_RESPONSE)

	// Rotate keys
	state.Session.Kesk = state.Session.NextKesk
	state.Session.Kepk = state.Session.NextKepk

	// set next round keys to nil
	state.Session.NextKesk = nil
	state.Session.NextKepk = nil
}

// checks for potential errors
// sends message to server
func send(conn *websocket.Conn, request []byte, err error, code int) {
	checkError(err)
	err = conn.WriteMessage(websocket.BinaryMessage, request)
	checkError(err)

	log.WithFields(log.Fields{
		"code": code,
	}).Info("SENT: ")
}

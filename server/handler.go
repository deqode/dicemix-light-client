package server

import (
	"fmt"
	"log"
	"os"
	"time"

	"../commons"
	"../utils"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	ecdh "github.com/wsddn/go-ecdh"
)

// identifies response message from server
// and passes response to appropriate handle for further operations
func handleMessage(conn *websocket.Conn, message []byte, code uint32, state *utils.State) {
	switch code {
	case commons.S_JOIN_RESPONSE:
		// Response against request to join dicemix transaction
		response := &commons.RegisterResponse{}
		err := proto.Unmarshal(message, response)
		checkError(err)
		handleJoinResponse(response, state)
	case commons.S_START_DICEMIX:
		// Response to start DiceMix Run
		response := &commons.DiceMixResponse{}
		err := proto.Unmarshal(message, response)
		checkError(err)
		handleStartDicemix(conn, response, state)
	case commons.S_KEY_EXCHANGE:
		// Response against request for KeyExchange
		response := &commons.DiceMixResponse{}
		err := proto.Unmarshal(message, response)
		checkError(err)
		handleKeyExchangeResponse(conn, response, state)
	case commons.S_EXP_DC_VECTOR:
		// contains roots of DC-Combined
		response := &commons.DCExpResponse{}
		err := proto.Unmarshal(message, response)
		checkError(err)
		handleDCExpResponse(conn, response, state)
	case commons.S_SIMPLE_DC_VECTOR:
		// conatins peers DC-SIMPLE-VECTOR's
		response := &commons.DiceMixResponse{}
		err := proto.Unmarshal(message, response)
		checkError(err)
		handleDCSimpleResponse(conn, response, state)
	case commons.S_TX_CONFIRMATION:
		// conatins peers DC-SIMPLE-VECTOR's
		response := &commons.DiceMixResponse{}
		err := proto.Unmarshal(message, response)
		checkError(err)
		handleConfirmationResponse(conn, response, state)
	}
}

// Response against request to join dicemix transaction
func handleJoinResponse(response *commons.RegisterResponse, state *utils.State) {
	if response.Err != "" {
		fmt.Fprintf(os.Stderr, "error: %v\n", response.Err)
		os.Exit(1)
	}
	// stores MyId provided by user
	state.MyID = response.Id

	fmt.Printf("\n%v\n", response.Message)
	fmt.Printf("MY ID - %v\n", state.MyID)
}

// Response to start DiceMix Run
func handleStartDicemix(conn *websocket.Conn, response *commons.DiceMixResponse, state *utils.State) {
	if response.Err != "" {
		fmt.Fprintf(os.Stderr, "error: %v\n", response.Err)
		os.Exit(1)
	}

	// initialize variables
	state.Peers = make([]utils.Peers, len(response.Peers)-1)
	set := make(map[int32]struct{}, len(response.Peers)-1)
	i := 0

	// store peers ID's
	// check for duplicate peer id's
	for _, peer := range response.Peers {
		if _, ok := set[peer.Id]; ok {
			log.Fatal("Duplicate peer IDs:", peer.Id)
			os.Exit(1)
		}
		set[peer.Id] = struct{}{}

		if peer.Id != state.MyID {
			state.Peers[i].ID = peer.Id
			i++
		}
	}

	// generates NIKE KeyPair for current run
	// mode = 0 to generate (my_kesk, my_kepk)
	iNike.GenerateKeys(state, 0)

	fmt.Printf("MY KEPK - %v\n", state.Kepk)

	// KeyExchange
	// broadcast our NIKE PublicKey with our peers
	ecdh := ecdh.NewCurve25519ECDH()
	keyExchangeRequest, err := proto.Marshal(&commons.KeyExchangeRequest{
		Code:      commons.C_KEY_EXCHANGE,
		Id:        state.MyID,
		PublicKey: ecdh.Marshal(state.Kepk),
		NumMsgs:   state.MyMsgCount,
		Timestamp: timestamp(),
	})

	// broadcast our PublicKey
	broadcast(conn, keyExchangeRequest, err)
}

// Response against request for KeyExchange
func handleKeyExchangeResponse(conn *websocket.Conn, response *commons.DiceMixResponse, state *utils.State) {
	if response.Err != "" {
		fmt.Fprintf(os.Stderr, "error: %v\n", response.Err)
		os.Exit(1)
	}

	// generate random 160 bit message
	for i := 0; i < int(state.MyMsgCount); i++ {
		state.MyMessages[i] = utils.GenerateMessage()
	}

	fmt.Printf("MY MESSAGE - %v\n\n", utils.Base58StringToBytes(state.MyMessages[0]))

	// copies peers info returned from server to local state.Peers
	// store peers PublicKey and NumMsgs
	filterPeers(state, response.Peers)

	// derive shared keys with peers
	iNike.DeriveSharedKeys(state)

	// generate DC Exponential Vector
	iDcNet.DeriveMyDCVector(state)

	// DC EXP
	// broadcast our DC-EXP vector with peers
	dcExpRequest, err := proto.Marshal(&commons.DCExpRequest{
		Code:        commons.C_EXP_DC_VECTOR,
		Id:          state.MyID,
		DCExpVector: state.MyDC,
		Timestamp:   timestamp(),
	})

	// broadcast our my_dc[]
	log.Printf("BR DC EXP")
	broadcast(conn, dcExpRequest, err)
}

// obtains roots and runs DC_SIMPLE
func handleDCExpResponse(conn *websocket.Conn, response *commons.DCExpResponse, state *utils.State) {
	if response.Err != "" {
		fmt.Fprintf(os.Stderr, "error: %v\n", response.Err)
		os.Exit(1)
	}

	// store roots (message hashes) calculated by server
	state.AllMsgHashes = response.Roots

	fmt.Printf("\nRECV: Roots - %v\n", state.AllMsgHashes)

	// run a SIMPLE DC NET
	iDcNet.RunDCSimple(state)

	if state.NextKepk == nil {
		// generates NIKE KeyPair for next run
		// mode = 1 to generate (my_next_kesk, my_next_kepk)
		iNike.GenerateKeys(state, 1)
	}

	// broadcast our DC SIMPLE Vector
	ecdh := ecdh.NewCurve25519ECDH()
	dcSimpleRequest, err := proto.Marshal(&commons.DCSimpleRequest{
		Code:           commons.C_SIMPLE_DC_VECTOR,
		Id:             state.MyID,
		DCSimpleVector: state.DCSimpleVector,
		MyOk:           state.MyOk,
		NextPublicKey:  ecdh.Marshal(state.NextKepk),
		Timestamp:      timestamp(),
	})

	broadcast(conn, dcSimpleRequest, err)
}

// handles other peers DC-SIMPLE-VECTORS
// resolves DC-NET
func handleDCSimpleResponse(conn *websocket.Conn, response *commons.DiceMixResponse, state *utils.State) {
	if response.Err != "" {
		fmt.Fprintf(os.Stderr, "error: %v\n", response.Err)
		os.Exit(1)
	}

	// copies peers info returned from server to local state.Peers
	// store other peers DC Simple Vectors
	filterPeers(state, response.Peers)

	// finally resolves DC Net Vectors to obtain messages
	// should contain all honest peers messages in absence of malicious peers
	iDcNet.ResolveDCNet(state)

	// Verify that every peer agrees to proceed
	ok := iDcNet.VerifyProceed(state)

	fmt.Printf("\nAgree to Proceed? = %v\n", ok)

	// broadcast our Confirmation
	confirmationRequest, err := proto.Marshal(&commons.ConfirmationRequest{
		Code:      commons.C_TX_CONFIRMATION,
		Id:        state.MyID,
		Confirm:   ok,
		Messages:  state.AllMessages,
		Timestamp: timestamp(),
	})

	broadcast(conn, confirmationRequest, err)
}

// handles other peers Confirmations
func handleConfirmationResponse(conn *websocket.Conn, response *commons.DiceMixResponse, state *utils.State) {
	if response.Err != "" {
		fmt.Fprintf(os.Stderr, "error: %v\n", response.Err)
		os.Exit(1)
	}

	success := state.MyOk

	// copies peers info returned from server to local state.Peers
	// store other peers Confirmations
	filterPeers(state, response.Peers)

	for i := 0; i < len(state.Peers); i++ {
		success = success && state.Peers[i].Confirm
	}

	// if every peer agrees to continue
	if success {
		fmt.Printf("\n\nTransaction successfull. All peers agreed.\n\n")
		conn.Close()
	} else {
		// else move to blame stage
		fmt.Printf("\n\nError occured. Need to find the culprit\n\n")
	}
}

// send request to server
func broadcast(conn *websocket.Conn, request []byte, err error) {
	checkError(err)
	err = conn.WriteMessage(websocket.BinaryMessage, request)
	checkError(err)
}

// to identify time of occurence of an event
// returns current timestamp
// example - 2018-08-07 12:04:46.456601867 +0000 UTC m=+0.000753626
func timestamp() string {
	return time.Now().String()
}

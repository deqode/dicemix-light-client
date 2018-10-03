package server

import (
	"time"

	"github.com/techracers-blockchain/dicemix-light-client/ecdsa"
	"github.com/techracers-blockchain/dicemix-light-client/messages"
	"github.com/techracers-blockchain/dicemix-light-client/utils"

	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

// copies peers info returned from server to local state.Peers
func filterPeers(state *utils.State, peers []*messages.PeersInfo) {
	// insanity check
	// if server sends more peers than actually involved in run
	// +1 represents peer himself, as server broadcast all clients info including his
	if len(state.Peers)+1 < len(peers) {
		log.Fatal("Error: obtained more peers from that we started. Expected - ", len(state.Peers), ", Obtained - ", len(peers))
	}

	// stores current peers info
	var peerIDs = make(map[int32]utils.Peers)

	// map peersInfo to their ids
	for _, peer := range state.Peers {
		peerIDs[peer.ID] = peer
	}

	// clears current available info of peers
	state.Peers = make([]utils.Peers, 0)

	for _, peer := range peers {
		// insanity check - to exclude clients sent by server which is not one of our peers
		// works for the case if server malfunctions
		if _, ok := peerIDs[peer.Id]; !ok {
			continue
		}

		// store peers updated info
		var tempPeer utils.Peers

		tempPeer.ID = peer.Id
		tempPeer.PubKey = peer.PublicKey
		tempPeer.NumMsgs = peer.NumMsgs
		tempPeer.SharedKey = peerIDs[peer.Id].SharedKey
		tempPeer.Dicemix = peerIDs[peer.Id].Dicemix
		tempPeer.DCSimpleVector = peer.DCSimpleVector
		tempPeer.Ok = peer.OK
		tempPeer.Confirmation = peer.Confirmation

		// add peer info to our peers
		state.Peers = append(state.Peers, tempPeer)
	}
}

// generates a RequestHeader proto
func requestHeader(code uint32, sessionID uint64, id int32) *messages.RequestHeader {
	return &messages.RequestHeader{
		Code:      code,
		SessionId: sessionID,
		Id:        id,
		Timestamp: timestamp(),
	}
}

// signs the message with privateKey and returns a Marshalled SignedRequest proto.
// It will panic if len(privateKey) is not PrivateKeySize.
func generateSignedRequest(privateKey, message []byte) ([]byte, error) {
	ecdsa := ecdsa.NewCurveECDSA()

	return proto.Marshal(&messages.SignedRequest{
		RequestData: message,
		Signature:   ecdsa.Sign(privateKey, message),
	})
}

// checks for any potential errors
// exists program if one found
func checkError(err error) {
	if err != nil {
		log.Fatalf("Error - %v", err)
	}
}

// to identify time of occurence of an event
// returns current timestamp
// example - 2018-08-07 12:04:46.456601867 +0000 UTC m=+0.000753626
func timestamp() string {
	return time.Now().String()
}

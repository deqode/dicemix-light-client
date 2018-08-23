package server

import (
	"time"

	"../eddsa"
	"../messages"
	"../utils"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/copier"
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

	var peersInfo []utils.Peers
	copier.Copy(&peersInfo, &state.Peers)
	state.Peers = make([]utils.Peers, 0)

	for i := 0; i < len(peers); i++ {
		for j := 0; j < len(peersInfo); j++ {
			if peers[i].Id != peersInfo[j].ID {
				continue
			}

			var tempPeer utils.Peers

			tempPeer.ID = peers[i].Id
			tempPeer.PubKey = peers[i].PublicKey
			tempPeer.NextPubKey = peers[i].NextPublicKey
			tempPeer.NumMsgs = peers[i].NumMsgs
			tempPeer.SharedKey = peersInfo[j].SharedKey
			tempPeer.Dicemix = peersInfo[j].Dicemix
			tempPeer.DCSimpleVector = peers[i].DCSimpleVector
			tempPeer.Ok = peers[i].OK
			tempPeer.Confirmation = peers[i].Confirmation

			state.Peers = append(state.Peers, tempPeer)
			break
		}
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
	edDSA := eddsa.NewCurveED25519()

	return proto.Marshal(&messages.SignedRequest{
		RequestData: message,
		Signature:   edDSA.Sign(privateKey, message),
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

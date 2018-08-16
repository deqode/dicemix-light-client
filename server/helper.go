package server

import (
	"../commons"
	"../utils"
	"github.com/jinzhu/copier"
)

// copies peers info returned from server to local state.Peers
func filterPeers(state *utils.State, peers []*commons.PeersInfo) {
	var peersInfo []utils.Peers
	copier.Copy(&peersInfo, &state.Peers)
	state.Peers = make([]utils.Peers, 0)

	for i := 0; i < len(peers); i++ {
		for j := 0; j < len(peersInfo); j++ {
			if peers[i].Id == peersInfo[j].ID {
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
}

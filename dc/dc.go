package dc

import (
	"dicemix-client/utils"
)

// DC - The main interface DC_NET.
type DC interface {
	DeriveMyDCVector(*utils.State)
	RunDCSimple(*utils.State)
	VerifyProceed(state *utils.State) bool
}

package dc

import (
	"../utils"
)

// DC - The main interface DC_NET.
type DC interface {
	DeriveMyDCVector(*utils.State)
	RunDCSimple(*utils.State)
	ResolveDCNet(*utils.State)
}

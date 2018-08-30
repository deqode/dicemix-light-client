package server

import (
	"github.com/manjeet-thadani/dicemix-client/utils"
)

// Server - The main interface to enable connection with server.
type Server interface {
	Register(*utils.State)
}

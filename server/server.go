package server

import (
	"github.com/techracers-blockchain/dicemix-light-client/utils"
)

// Server - The main interface to enable connection with server.
type Server interface {
	Register(*utils.State)
}

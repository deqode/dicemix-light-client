package server

import (
	"dicemix_client/utils"
)

// Server - The main interface to enable connection with server.
type Server interface {
	Register(*utils.State)
}

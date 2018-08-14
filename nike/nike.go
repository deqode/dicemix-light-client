package nike

import (
	"../utils"
)

// NIKE - The main interface for Non-interactive Key Exchange (NIKE).
type NIKE interface {
	GenerateKeys(*utils.State, int)
	KeyExchange(*utils.State)
	DeriveSharedKeys(*utils.State)
}

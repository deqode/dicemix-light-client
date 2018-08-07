package nike

import (
	"../utils"
)

// NIKE - The main interface for Non-interactive Key Exchange (NIKE).
type NIKE interface {
	GenerateKeys(*utils.State)
	KeyExchange(*utils.State)
	DeriveSharedKeys(*utils.State)
}

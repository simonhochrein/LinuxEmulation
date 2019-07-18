package interfaces

import (
	"github.com/unicorn-engine/unicorn/bindings/go/unicorn"
)

// Emulator - Main class
type Emulator interface {
	GetUnicorn() unicorn.Unicorn
	GetStackManager() StackManager
	GetTypeConversion() TypeConversion
	GetBRK() uint64
	SetBRK(brk uint64) uint64
}

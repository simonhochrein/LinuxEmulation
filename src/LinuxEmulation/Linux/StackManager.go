package linux

import (
	helpers "LinuxEmulation/Helpers"
	interfaces "LinuxEmulation/Interfaces"

	"github.com/unicorn-engine/unicorn/bindings/go/unicorn"
)

const (
	stackBase = 0xbf800000
	stackSize = 0x00800000
)

// StackManager for emulator
type StackManager struct {
	emulator interfaces.Emulator
}

// Push - push data onto the stack
func (s StackManager) Push(bytes []byte) uint64 {
	uc := s.emulator.GetUnicorn()
	rsp, _ := uc.RegRead(unicorn.X86_REG_RSP)
	rsp -= uint64(len(bytes))
	uc.RegWrite(unicorn.X86_REG_RSP, rsp)
	uc.MemWrite(rsp, bytes)
	return rsp
}

// func (s StackManager) Pop() {

// }

// Align - Align the stack pointer
func (s StackManager) Align() {
	uc := s.emulator.GetUnicorn()
	rsp, _ := uc.RegRead(unicorn.X86_REG_RSP)
	rsp = rsp &^ 15
	uc.RegWrite(unicorn.X86_REG_RSP, rsp)
}

// Map - Map initial stack memory to emulator instance + Auxv
func (s StackManager) Map() {
	uc := s.emulator.GetUnicorn()

	uc.MemMap(stackBase, stackSize)
	uc.RegWrite(unicorn.X86_REG_RSP, stackBase+stackSize)
	s.Push(helpers.EncodeUint64(0))
	argvp := s.Push(helpers.EncodeString("./a.out\x00"))
	s.Align()
	// AUXV
	s.Push(helpers.EncodeUint64(0))
	s.Push(helpers.EncodeUint64(0))
	//
	s.Push(helpers.EncodeUint64(0))
	s.Push(helpers.EncodeUint64(0))
	s.Push(helpers.EncodeUint64(argvp))
	s.Push(helpers.EncodeUint64(1))
}

// NewStackManager - create a new instance of the stack manager
func NewStackManager(emulator interfaces.Emulator) *StackManager {
	return &StackManager{emulator: emulator}
}

package linuxemulation

import (
	debug "LinuxEmulation/Debug"
	interfaces "LinuxEmulation/Interfaces"
	linux "LinuxEmulation/Linux"
	"debug/elf"
	"fmt"

	"github.com/unicorn-engine/unicorn/bindings/go/unicorn"
)

// Emulator - Main class
type Emulator struct {
	internet       interfaces.Internet
	filesystem     interfaces.FileSystem
	syscalls       *linux.SystemCalls
	unicorn        unicorn.Unicorn
	stackManager   *linux.StackManager
	typeConversion *linux.TypeConversion
	brk            uint64
	capstone       *debug.Capstone
}

//Run - Run a program
func (e *Emulator) Run(program *elf.File) {
	linux.LoadELF(program, e)
	e.stackManager.Map()
	err := e.unicorn.Start(program.Entry, 0)
	if err != nil {
		panic(err)
	}
}

func (e *Emulator) hooks() {
	e.unicorn.HookAdd(unicorn.HOOK_INSN, func(uc unicorn.Unicorn) {
		e.syscalls.HandleCall()
	}, 1, 0, unicorn.X86_INS_SYSCALL)
	e.unicorn.HookAdd(unicorn.HOOK_MEM_UNMAPPED, func(uc unicorn.Unicorn, access int, address uint64, size int, value int64) bool {
		fmt.Printf("%x: %x\n", address, size)
		return false
	}, 1, 0)
	e.unicorn.HookAdd(unicorn.HOOK_CODE, func(uc unicorn.Unicorn, addr uint64, size uint32) {
		bytes, _ := uc.MemRead(addr, uint64(size))
		val, _ := uc.RegRead(unicorn.X86_REG_RAX)
		e.capstone.Disassemble(bytes, addr, uint64(size), val)
	}, 1, 0)
}

// GetUnicorn - Return the current instance of Unicorn Engine
func (e *Emulator) GetUnicorn() unicorn.Unicorn {
	return e.unicorn
}

// GetStackManager - Return the current instance of the stack manager
func (e *Emulator) GetStackManager() interfaces.StackManager {
	return e.stackManager
}

// GetTypeConversion - Return the current instance of the type converter
func (e *Emulator) GetTypeConversion() interfaces.TypeConversion {
	return e.typeConversion
}

// GetBRK - Return the current brk
func (e *Emulator) GetBRK() uint64 {
	return e.brk
}

// SetBRK - Set brk for the current instance
func (e *Emulator) SetBRK(brk uint64) uint64 {
	// Do allocation stuff here
	e.brk = brk
	return e.brk
}

// NewEmulator - Create a new instance of the emulator
func NewEmulator(internet interfaces.Internet, filesystem interfaces.FileSystem) *Emulator {
	uc, _ := unicorn.NewUnicorn(unicorn.ARCH_X86, unicorn.MODE_64)
	emu := &Emulator{internet: internet, filesystem: filesystem, unicorn: uc, capstone: debug.NewCapstone()}
	emu.stackManager = linux.NewStackManager(emu)
	emu.syscalls = linux.NewSystemCalls(emu)
	emu.typeConversion = linux.NewTypeConversion(emu)
	emu.hooks()
	return emu
}

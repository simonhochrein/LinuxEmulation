package linux

import (
	helpers "LinuxEmulation/Helpers"
	interfaces "LinuxEmulation/Interfaces"
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/unicorn-engine/unicorn/bindings/go/unicorn"
)

// SystemCalls - Handler for syscall
type SystemCalls struct {
	emulator interfaces.Emulator
}

type registers = []uint64

// HandleCall - Handle System Calls
func (s SystemCalls) HandleCall() {
	uc := s.emulator.GetUnicorn()
	regs, _ := uc.RegReadBatch([]int{unicorn.X86_REG_RAX, unicorn.X86_REG_RDI, unicorn.X86_REG_RSI, unicorn.X86_REG_RDX, unicorn.X86_REG_R10, unicorn.X86_REG_R8, unicorn.X86_REG_R9})
	switch regs[0] {
	case 2:
		s.open(regs)
		break
	case 12:
		s.brk(regs)
		break
	case 20:
		s.writev(regs)
		break
	case 63:
		s.uname(regs)
		break
	default:
		panic(fmt.Sprintf("Unknown syscall: %d", regs[0]))
	}
}

func (s SystemCalls) brk(regs registers) {
	addr := s.emulator.SetBRK(regs[1])
	s.ret(addr)
}

func (s SystemCalls) open(regs registers) {
	typeConversion := s.emulator.GetTypeConversion()
	println(typeConversion.ReadCString(regs[1]))
	//TODO: this is a hack, but it does the trick for now
	s.ret(0)
}

type iovec struct {
	IovBase uint64 /* Starting address */
	IovLen  uint64 /* Number of bytes to transfer */
}

func (s SystemCalls) writev(regs registers) {
	uc := s.emulator.GetUnicorn()
	fmt.Printf("FD: %d, length: %d\n", regs[1], regs[3])
	for i := uint64(0); i < regs[3]; i++ {
		v := iovec{}
		offset := regs[2] + i*16 /*Size of two uint64s*/
		b, e := uc.MemRead(offset, 16)
		if e != nil {
			panic(e)
		}
		binary.Read(bytes.NewReader(b), binary.LittleEndian, &v)
		str, _ := uc.MemRead(v.IovBase, v.IovLen)
		println(string(str))
	}
	s.ret(0)
}

type utsname struct {
	sysname  [65]byte
	nodename [65]byte
	release  [65]byte
	version  [65]byte
	machine  [65]byte
}

func (s SystemCalls) uname(regs registers) {
	ret := utsname{}
	copy(ret.release[:], []byte("3.13.0-24-generic\x00"))
	s.emulator.GetUnicorn().MemWrite(regs[1], helpers.PackStruct(ret))
	s.ret(0)
}

func (s SystemCalls) ret(val uint64) {
	s.emulator.GetUnicorn().RegWrite(unicorn.X86_REG_RAX, val)
}

// NewSystemCalls - Create an instance of the system call handler
func NewSystemCalls(emulator interfaces.Emulator) *SystemCalls {
	return &SystemCalls{emulator: emulator}
}

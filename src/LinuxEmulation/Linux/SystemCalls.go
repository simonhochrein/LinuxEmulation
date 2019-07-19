package linux

import (
	helpers "LinuxEmulation/Helpers"
	interfaces "LinuxEmulation/Interfaces"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"

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
	fmt.Printf("SYSCALL %d\n", regs[0])
	switch regs[0] {
	case 2:
		s.open(regs)
		break
	case 9:
		s.mmap(regs)
		break
	case 12:
		s.brk(regs)
		break
	case 16:
		//TODO: IOCTL
		s.ret(0)
		break
	case 20:
		s.writev(regs)
		break
	case 63:
		s.uname(regs)
		break

	case 158:
		s.archPrctl(regs)
		break
	default:
		panic(fmt.Sprintf("Unknown syscall: %d", regs[0]))
	}
}

func (s SystemCalls) brk(regs registers) {
	addr := regs[1]
	oldBrk := s.emulator.GetBRK()
	uc := s.emulator.GetUnicorn()

	if addr > 0 && addr >= oldBrk {
		fmt.Printf("%x : %x\n", oldBrk, addr)
		size := addr - oldBrk
		uc.MemMap(oldBrk, helpers.PageEnd(size))
		s.ret(addr)
		return
		// panic("")
	}
	s.ret(oldBrk)
}

func (s SystemCalls) open(regs registers) {
	typeConversion := s.emulator.GetTypeConversion()
	println(typeConversion.ReadCString(regs[1]))
	//TODO: this is a hack, but it does the trick for now
	s.ret(0)
}

func (s SystemCalls) mmap(regs registers) {
	uc := s.emulator.GetUnicorn()
	addr, size, prot, flags, fd, offset := regs[1], regs[2], regs[3], regs[4], int64(regs[5]), regs[6]
	fmt.Printf("mmap(0x%x, 0x%x, 0x%x, 0x%x, %d, 0x%x)\n", addr, size, prot, flags, fd, offset)

	if fd > 0 {

	}
	uc.MemMap(addr, size)
	// s.ret(addr)
	os.Exit(0)
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

const (
	ARCH_SET_GS = 0x1001
	ARCH_SET_FS = 0x1002
	ARCH_GET_FS = 0x1003
	ARCH_GET_GS = 0x1004
)

func (s SystemCalls) archPrctl(regs registers) {
	uc := s.emulator.GetUnicorn()
	fsmsr := uint64(0xC0000100)
	gsmsr := uint64(0xC0000101)

	addr := regs[2]
	// TODO: make SET check for valid mapped memory
	switch regs[1] {
	case ARCH_SET_FS:
		uc.RegWriteX86Msr(fsmsr, addr)
	case ARCH_SET_GS:
		uc.RegWriteX86Msr(gsmsr, addr)
	case ARCH_GET_FS:
		val, _ := uc.RegReadX86Msr(fsmsr)
		buf := helpers.EncodeUint64(val)
		uc.MemWrite(addr, buf)
	case ARCH_GET_GS:
		val, _ := uc.RegReadX86Msr(gsmsr)
		buf := helpers.EncodeUint64(val)
		uc.MemWrite(addr, buf)
	}
}

func (s SystemCalls) ret(val uint64) {
	s.emulator.GetUnicorn().RegWrite(unicorn.X86_REG_RAX, val)
}

// NewSystemCalls - Create an instance of the system call handler
func NewSystemCalls(emulator interfaces.Emulator) *SystemCalls {
	return &SystemCalls{emulator: emulator}
}

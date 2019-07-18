package syscalls

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"utils"

	"github.com/unicorn-engine/unicorn/bindings/go/unicorn"
)

var BRK uint64

func ret(uc unicorn.Unicorn, n int) {
	uc.RegWrite(unicorn.X86_REG_RAX, uint64(n))
}

func ret64(uc unicorn.Unicorn, n uint64) {
	uc.RegWrite(unicorn.X86_REG_RAX, n)
}

func dumpBytes(n uint64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, n)
	return buf
}

func dumpStruct(data interface{}) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, data)
	return buf.Bytes()
}

func readCStr(uc unicorn.Unicorn, p uint64) {
	var current byte
	var str string
	p2 := p
	for current != 0 || p2 == p {
		bytes, err := uc.MemRead(p2, 1)
		if err != nil {
			panic(err)
		}
		str = str + string(current)
		current = bytes[0]
		p2++
	}
	println(str)
}

func HandleSyscall(uc unicorn.Unicorn) {
	rax, _ := uc.RegRead(unicorn.X86_REG_RAX)
	switch rax {
	case 2:
		Open(uc)
		break
	case 12:
		Brk(uc)
		break
	case 20:
		WriteV(uc)
		break
	case 63:
		Uname(uc)
		break
	default:
		panic(fmt.Sprintf("Unimplemented Syscall: %v", rax))
	}
}

type oldUtsname struct {
	sysname  [65]byte
	nodename [65]byte
	release  [65]byte
	version  [65]byte
	machine  [65]byte
}

func Open(uc unicorn.Unicorn) {
	rdi, _ := uc.RegRead(unicorn.X86_REG_RDI)
	readCStr(uc, rdi)
	ret(uc, 0)
}

func Brk(uc unicorn.Unicorn) {
	rdi, _ := uc.RegRead(unicorn.X86_REG_RDI)
	if rdi > 0 && rdi >= BRK {
		// prot := elf.PF_R | elf.PF_W
		size := rdi - BRK
		println(uint64(rdi), ":", BRK, ":", size)
		if size > 0 {

			alignedStart := utils.PageStart(rdi)
			alignedSize := utils.PageEnd((rdi - alignedStart) + size)

			err := uc.MemMap(alignedStart, alignedSize)
			if err != nil {
				panic(err)
			}
			println("EXPAND")
		}
		BRK = rdi
	}
	// readCStr(uc, rdi)
	ret64(uc, BRK)
}

func Uname(uc unicorn.Unicorn) {
	rdi, _ := uc.RegRead(unicorn.X86_REG_RDI)
	var version [65]byte
	copy(version[:], []byte("3.13.0-24-generic\x00"))
	// println(string(version[:]))
	var uname oldUtsname = oldUtsname{release: version}
	uc.MemWrite(rdi, dumpStruct(uname))
	ret(uc, 0)
}

type iovec struct {
	iov_base uint64
	iov_len  uint64
}

func WriteV(uc unicorn.Unicorn) {
	rdi, _ := uc.RegRead(unicorn.X86_REG_RDI)
	rsi, _ := uc.RegRead(unicorn.X86_REG_RSI)
	rdx, _ := uc.RegRead(unicorn.X86_REG_RDX)
	println(string(rdi))
	for i := uint64(0); i < rdx; i++ {
		bytes, _ := uc.MemRead(rsi+i*16, 16)
		v := iovec{iov_base: binary.LittleEndian.Uint64(bytes[:8]), iov_len: binary.LittleEndian.Uint64(bytes[8:])}
		b, _ := uc.MemRead(v.iov_base, v.iov_len)
		println(string(b))
	}
}

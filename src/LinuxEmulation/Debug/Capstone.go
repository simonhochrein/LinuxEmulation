package debug

import (
	"fmt"

	"github.com/lunixbochs/capstr"
)

// Capstone Disassembler
type Capstone struct {
	engine *capstr.Engine
}

func (c *Capstone) Disassemble(bytes []byte, address uint64, size uint64, rax uint64) {
	instrs, _ := c.engine.Dis(bytes, address, size)
	for _, instr := range instrs {
		fmt.Printf("%x: %s %s :: %x\n", instr.Addr(), instr.Mnemonic(), instr.OpStr(), rax)
	}
}

// NewCapstone - create instance of capstone disassembler
func NewCapstone() *Capstone {
	engine, _ := capstr.New(capstr.ARCH_X86, capstr.MODE_64)
	return &Capstone{engine}
}

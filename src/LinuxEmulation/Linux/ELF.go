package linux

import (
	helpers "LinuxEmulation/Helpers"
	interfaces "LinuxEmulation/Interfaces"
	"debug/elf"
	"fmt"
)

// LoadELF - loads an elf file into the emulator
func LoadELF(file *elf.File, emulator interfaces.Emulator) {
	uc := emulator.GetUnicorn()
	brk := uint64(0)
	for _, programSection := range file.Progs {
		if programSection.Type == elf.PT_LOAD {
			data := make([]byte, programSection.Filesz)

			programSection.Open().Read(data)

			alignedStart := helpers.PageStart(programSection.Vaddr)
			alignedSize := helpers.PageEnd((programSection.Vaddr - alignedStart) + programSection.Memsz)

			fmt.Printf("Map 0x%x with length 0x%x\n", alignedStart, alignedSize)

			uc.MemMap(alignedStart, alignedSize)
			uc.MemWrite(programSection.Vaddr, data)
		}
		if programSection.Flags&elf.PF_W != 0 {
			addr := programSection.Vaddr + programSection.Memsz
			if addr > emulator.GetBRK() {
				brk = addr
			}
		}
	}
	if brk > 0 {
		emulator.SetBRK(helpers.PageEnd(brk))
	}
}

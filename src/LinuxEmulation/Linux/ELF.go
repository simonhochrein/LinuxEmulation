package linux

import (
	helpers "LinuxEmulation/Helpers"
	interfaces "LinuxEmulation/Interfaces"
	"debug/elf"
)

// LoadELF - loads an elf file into the emulator
func LoadELF(file *elf.File, emulator interfaces.Emulator) {
	uc := emulator.GetUnicorn()
	for _, programSection := range file.Progs {
		if programSection.Type == elf.PT_LOAD {
			data := make([]byte, programSection.Filesz)
			programSection.Open().Read(data)

			alignedStart := helpers.PageStart(programSection.Vaddr)
			alignedSize := helpers.PageEnd((programSection.Vaddr - alignedStart) + programSection.Memsz)

			uc.MemMap(alignedStart, alignedSize)
			uc.MemWrite(programSection.Vaddr, data)
		}
	}
}

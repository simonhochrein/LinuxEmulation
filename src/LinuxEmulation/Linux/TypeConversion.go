package linux

import interfaces "LinuxEmulation/Interfaces"

// TypeConversion - Convert between linux types and go types
type TypeConversion struct {
	emulator interfaces.Emulator
}

// ReadCString - Read a C string from memory
func (t TypeConversion) ReadCString(addr uint64) string {
	uc := t.emulator.GetUnicorn()
	var c byte
	str := ""
	index := uint64(0)

	for c != '\x00' || index == 0 {
		b, e := uc.MemRead(addr+index, 1)
		if e != nil {
			panic("Segmentation Fault")
		}
		c = b[0]
		str += string(c)
		index++
	}
	return str
}

// NewTypeConversion - Create a new instance of the type conversion class
func NewTypeConversion(emulator interfaces.Emulator) *TypeConversion {
	return &TypeConversion{emulator: emulator}
}

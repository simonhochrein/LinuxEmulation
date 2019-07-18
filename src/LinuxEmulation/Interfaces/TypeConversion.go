package interfaces

// TypeConversion - Type definition for TypeConversion class
type TypeConversion interface {
	ReadCString(addr uint64) string
}

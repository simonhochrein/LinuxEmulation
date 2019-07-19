package helpers

// PageStart - Round page down
func PageStart(n uint64) uint64 {
	return n &^ (4*1024 - 1)
}

// PageEnd - Round page up
func PageEnd(n uint64) uint64 {
	return (n + 4*1024 - 1) &^ (4*1024 - 1)
}

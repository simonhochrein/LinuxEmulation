package utils

func PageStart(n uint64) uint64 {
	return n &^ (8*1024 - 1)
}

func PageEnd(n uint64) uint64 {
	return (n + 8*1024 - 1) &^ (8*1024 - 1)
}

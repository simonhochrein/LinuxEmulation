package interfaces

// StackManager - Interface for the class defined in Linux/Stack.go
type StackManager interface {
	Push(bytes []byte) uint64
}

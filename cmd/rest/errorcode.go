package main

const (
	ErrUnknownNode int = iota
	ErrConnNodeFailed
	ErrAddNodeFailed
	ErrNoNodeIsConnected
	ErrGetBlockFailed
	ErrGetTransactionFailed
	ErrInvalidParameter
)

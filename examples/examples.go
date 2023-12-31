package examples

// MatrixType represents mobile phone's LCD matrix type.
//
//go:generate iotas -type=MatrixType
type MatrixType uint8

const (
	Unknown MatrixType = iota
	OLED
	AMOLED
	TFT
)

type ExtraType int

const (
	One ExtraType = iota
	Two
	Three
)

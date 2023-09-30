package testpkg

type TestType uint8

const (
	SomeConst1 = "1"
)

const SomeConst2 = "2"

const (
	TestTypeX TestType = 0
	TestTypeY TestType = iota
)

type TestType2 uint8

const (
	TestType2X TestType2 = iota
)

const (
	TestType2Y TestType2 = iota
)

type TestType3 uint8

const (
	TestType3X TestType3 = iota
	TestType3Y           = TestType3X
)

type TestType4 uint8

const (
	TestType4X TestType4 = iota
	TestType4Y TestType4 = iota
)

package rand

type RandType int

const (
	ChaCha8 RandType = 1
	PCG     RandType = 2
	LCG     RandType = 3
	LCG2    RandType = 4

	Zipf        RandType = 10
	Exponential RandType = 11
	Normal      RandType = 11

	Default RandType = ChaCha8
)

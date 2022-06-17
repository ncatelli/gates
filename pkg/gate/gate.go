package gate

type IO bool

type Gate interface {
	InputCnt() uint
	InputPaths() []string
	CurrentTick() uint
	Tick(tick uint, inputs []IO)
}

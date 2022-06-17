package gate

type Not struct {
	tick uint
}

func (not *Not) Inputs() uint {
	return 1
}

func (not *Not) InputPaths() []string {
	return []string{"a"}
}

func (not *Not) CurrentTick() uint {
	return not.tick
}

func (not *Not) Tick(tick uint, inputs []IO) {
}

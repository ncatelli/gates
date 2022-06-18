package gate

import (
	"fmt"
	"unicode"
)

type IO bool

type Gate interface {
	Inputs() uint
	Compute(tick uint, input []IO) (IO, error)
}

type PendingInput struct {
	received bool
	state    IO
}

type TickState struct {
	inputs []PendingInput
}

func NewTickState(expectedInputs uint) *TickState {
	inputs := make([]PendingInput, expectedInputs)

	for i := 0; i < int(expectedInputs); i++ {
		inputs[i] = PendingInput{
			received: false,
			state:    false,
		}
	}

	return &TickState{
		inputs: inputs,
	}
}

// AllInputsReceived returns true if every pending input has been marked received.
func (ts *TickState) AllInputsReceived() bool {
	// check if all inputs have been received, if not return early.
	for _, input := range ts.inputs {
		if !input.received {
			return false
		}
	}

	return true
}

// ReturnInputsIfReady returns a slice of all IO state in order if called after
// all have been received. Otherwise an error is returned.
func (ts *TickState) ReturnInputsIfReady() ([]IO, error) {
	if !ts.AllInputsReceived() {
		return nil, fmt.Errorf("input is still pending")
	}

	inputs := make([]IO, 0, len(ts.inputs))
	for _, pending := range ts.inputs {
		inputs = append(inputs, pending.state)
	}

	return inputs, nil
}

type GenericGate struct {
	ticks          map[uint]*TickState
	expectedInputs uint
	gate           Gate
}

func NewGenericGate(g Gate) *GenericGate {
	return &GenericGate{
		ticks:          make(map[uint]*TickState),
		expectedInputs: g.Inputs(),
		gate:           g,
	}
}

func (hg *GenericGate) Inputs() uint {
	return hg.expectedInputs
}

func (hg *GenericGate) Compute(tick uint, inputs []IO) (IO, error) {
	if len(inputs) != int(hg.Inputs()) {
		return false, fmt.Errorf("input does not match expected length of %d", hg.Inputs())
	}

	return hg.gate.Compute(tick, inputs)

}

func (hg *GenericGate) ReceiveInput(tick uint, input rune, state IO) (*TickState, error) {
	// calculate the offset of the input from the rune (path)
	inputOffset, err := runeToNormalizedOffset(input)
	if err != nil {
		return nil, err
	}

	// verify that the input is within range
	if (inputOffset + 1) > hg.expectedInputs {
		return nil, fmt.Errorf("offset %v(%d) greater than max %d", input, inputOffset, hg.expectedInputs)
	}

	ts, prs := hg.ticks[tick]
	if !prs {
		ts = NewTickState(hg.expectedInputs)
		hg.ticks[tick] = ts
	}

	// error if inputs are clobbered.
	if ts.inputs[inputOffset].received {
		return nil, fmt.Errorf("input (%v) for tick %d already set", input, tick)
	}

	ts.inputs[inputOffset] = PendingInput{
		received: true,
		state:    state,
	}

	// check if all inputs have been received, if not return early.
	for _, input := range ts.inputs {
		if !input.received {
			return nil, nil
		}
	}

	// lastly if all inputs are ready, return the tick state, representing
	// completed tick.
	return ts, nil
}

func runeToNormalizedOffset(r rune) (uint, error) {
	min := uint('a')
	max := uint('z')
	normalized := unicode.ToLower(r)
	runeAsInt := uint(normalized)

	if runeAsInt <= max || runeAsInt >= min {
		offset := runeAsInt - min

		return offset, nil
	}

	return 0, fmt.Errorf("value out of range: must be between a-z, got %v", r)

}

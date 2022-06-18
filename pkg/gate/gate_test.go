package gate

import (
	"fmt"
	"testing"
)

type MockGate struct {
	inputs uint
}

func (mg *MockGate) Inputs() uint {
	return mg.inputs
}

// Compute echos the current value
func (mg *MockGate) Compute(tick uint, inputs []IO) (IO, error) {
	if len(inputs) != int(mg.Inputs()) {
		return false, fmt.Errorf("input does not match expected length of %d", mg.Inputs())
	}

	input := inputs[0]
	output := input

	return output, nil
}

func TestGateShouldReceiveInputShould(t *testing.T) {
	t.Run("return the TypeState when all valid inputs are received", func(t *testing.T) {
		gate := NewGenericGate(&MockGate{
			inputs: 1,
		})

		ts, err := gate.ReceiveInput(0, 'a', true)
		if err != nil {
			t.Error(err)
		} else if ts == nil {
			t.Error("ts should not be nil")
		}
	})

	t.Run("return all nil if all inputs aren't satisfied", func(t *testing.T) {
		gate := NewGenericGate(&MockGate{
			inputs: 2,
		})

		ts, err := gate.ReceiveInput(0, 'a', true)
		if err != nil {
			t.Error(err)
		} else if ts != nil {
			t.Error("ts should be nil")
		}
	})

	t.Run("error if inputs are clobbered", func(t *testing.T) {
		gate := NewGenericGate(&MockGate{
			inputs: 2,
		})

		// first input should pass with no error.
		_, err := gate.ReceiveInput(0, 'a', true)
		if err != nil {
			t.Error(err)
		}

		// second should clobber, triggering an error.
		_, err = gate.ReceiveInput(0, 'a', true)
		if err == nil {
			t.Error(err)
		}
	})
}

package gate

import (
	"fmt"
	"testing"

	"github.com/ncatelli/gates/pkg/models"
)

type noopOutputter struct{}

func (no *noopOutputter) Output(tick uint, state models.IO) error {
	return nil
}

type MockGate struct {
	inputs uint
}

func (mg *MockGate) Inputs() uint {
	return mg.inputs
}

// Compute echos the current value
func (mg *MockGate) Compute(tick uint, inputs []models.IO) (models.IO, error) {
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
		}, &noopOutputter{})

		ts, err := gate.ReceiveInput(0, 'a', true)
		if err != nil {
			t.Error(err)
		} else if ts == nil {
			t.Error("ts should not be nil")
		}
	})

	t.Run("return ts that flags non-ready if all inputs aren't satisfied", func(t *testing.T) {
		gate := NewGenericGate(&MockGate{
			inputs: 2,
		}, &noopOutputter{})

		ts, err := gate.ReceiveInput(0, 'a', true)
		if err != nil {
			t.Error(err)
		} else if ts.AllInputsReceived() {
			t.Error("ts should not be ready")
		}
	})

	t.Run("error if inputs are clobbered", func(t *testing.T) {
		gate := NewGenericGate(&MockGate{
			inputs: 2,
		}, &noopOutputter{})

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

package gate

import (
	"testing"

	"github.com/ncatelli/gates/pkg/models"
)

func TestOrShouldCompute(t *testing.T) {
	t.Run("the correct value for truth table", func(t *testing.T) {
		var tick uint = 0
		truthTable := [][3]models.IO{
			{false, false, false},
			{true, false, true},
			{false, true, true},
			{true, true, true},
		}

		for _, prop := range truthTable {
			gate := NewGenericGate(&Or{}, &noopOutputter{})

			a_in := prop[0]
			b_in := prop[1]
			expected_output := prop[2]

			_, err := gate.ReceiveInput(tick, 'a', a_in)
			if err != nil {
				t.Error(err)
			}
			ts, err := gate.ReceiveInput(tick, 'b', b_in)
			if err != nil {
				t.Error(err)
			} else if ts == nil {
				t.Error("ts should not be nil")
			}

			inputs, err := ts.ReturnInputsIfReady()
			if err != nil {
				t.Error(err)
			}

			output, err := gate.Compute(tick, inputs)
			if err != nil {
				t.Error(err)
			} else if output != expected_output {
				t.Errorf("wrong output for or gate: wanted %v got %v", expected_output, output)
			}
		}
	})

}

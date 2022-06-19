package gate

import (
	"testing"

	"github.com/ncatelli/gates/pkg/models"
)

func TestNotShouldCompute(t *testing.T) {
	t.Run("the inverse of a its input", func(t *testing.T) {
		var tick uint = 0
		truthTable := [][2]models.IO{
			{false, true},
			{true, false},
		}

		for _, prop := range truthTable {
			gate := NewGateService(&Not{}, &noopOutputter{})
			a_in := prop[1]
			expected_output := prop[1]

			ts, err := gate.ReceiveInput(tick, 'a', a_in)
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
			} else if output == expected_output {
				t.Errorf("wrong output for not gate: wanted %v got %v", false, output)
			}
		}
	})
}

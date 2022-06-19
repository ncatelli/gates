package gate

import (
	"fmt"

	"github.com/ncatelli/gates/pkg/models"
)

type Not struct{}

func (not *Not) Inputs() uint {
	return 1
}

func (not *Not) Compute(tick uint, inputs []models.IO) (models.IO, error) {
	if len(inputs) != int(not.Inputs()) {
		return false, fmt.Errorf("input does not match expected length of %d", not.Inputs())
	}

	input := inputs[0]
	output := !input

	return output, nil
}

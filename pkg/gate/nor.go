package gate

import (
	"fmt"

	"github.com/ncatelli/gates/pkg/models"
)

type Nor struct{}

func (nor *Nor) Inputs() uint {
	return 2
}

func (nor *Nor) Compute(tick uint, inputs []models.IO) (models.IO, error) {
	if len(inputs) != int(nor.Inputs()) {
		return false, fmt.Errorf("input does not match expected length of %d", nor.Inputs())
	}

	a := inputs[0]
	b := inputs[1]
	output := !(a || b)

	return output, nil
}

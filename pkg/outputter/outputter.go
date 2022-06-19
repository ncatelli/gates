package outputter

import "github.com/ncatelli/gates/pkg/models"

type Outputter interface {
	Output(uint, models.IO) error
}

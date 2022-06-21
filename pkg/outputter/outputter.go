package outputter

import "github.com/ncatelli/gates/pkg/models"

type OutputTy int

const (
	StdOut OutputTy = iota
	HTTP
)

type Outputter interface {
	Output(uint, models.IO) error
}

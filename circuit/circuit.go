package circuit

import (
	"github.com/consensys/gnark/frontend"
)

type Circuit struct {
}

// this is just a mock circuit, it will replaced by actual circuit logic
func (circuit *Circuit) Define(api frontend.API) error {
	return nil
}

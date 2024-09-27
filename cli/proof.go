package main

import (
	"log"

	zkprover "github.com/arka-labs/zk-prover"
	"github.com/arka-labs/zk-prover/circuit"
)

const (
	DIR_1 = "./keys"
)

func main() {
	var circuitInputs circuit.CircuitInputs

	circuitInputs.X = 1
	circuitInputs.Y = 2
	circuitInputs.Z = 3

	validityProof, err := zkprover.GetZKProof(circuitInputs, DIR_1)
	if err != nil {
		log.Fatal(err)
	}
	validityProof.Verify(DIR_1)
}

package main

import (
	"crypto/sha256"
	"log"

	zkprover "github.com/arka-labs/zk-prover"
	"github.com/arka-labs/zk-prover/circuit"
)

const (
	DIR_3 = "./keys"
)

func main() {
	var circuitInputs circuit.CircuitInputs

	data := []byte("good data")
	resizedData := circuit.FixedSizeBytes(data)

	digest := sha256.Sum256(resizedData[:])

	circuitInputs.Data = resizedData[:]
	circuitInputs.Hash = digest

	validityProof, err := zkprover.GetZKProof(circuitInputs, DIR_3)
	if err != nil {
		log.Fatal(err)
	}
	validityProof.Verify(DIR_3)
}

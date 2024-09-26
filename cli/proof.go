package main

import (
	"log"

	zkprover "github.com/arka-labs/zk-prover"
	"github.com/arka-labs/zk-prover/circuit"
	"github.com/consensys/gnark/backend/groth16"
)

func main() {
	zkProof, publicWitness, err := zkprover.GetZKProof(1, 2, 3)
	if err != nil {
		log.Fatal(err)
	}
	vk, err := circuit.GetVerifierKey()
	if err != nil {
		log.Fatal(err)
	}
	groth16.Verify(zkProof, vk, publicWitness)
}

package circuit

import (
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
)

const INPUTS_SIZE = 3

type CircuitInputs struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

type Circuit struct {
	X frontend.Variable `gnark:",public"`
	Y frontend.Variable `gnark:",public"`
	Z frontend.Variable `gnark:",public"`
}

// this is just a mock circuit, it will replaced by actual circuit logic
func (circuit *Circuit) Define(api frontend.API) error {
	return nil
}

func GenerateZKProof(cs constraint.ConstraintSystem, pk groth16.ProvingKey, customInputs any) (groth16.Proof, witness.Witness, error) {

	circuitInputs := customInputs.(CircuitInputs)

	// TODO: just change this
	// contruct assignment using dynamic veriables
	assignment := Circuit{
		X: circuitInputs.X,
		Y: circuitInputs.Y,
		Z: circuitInputs.Z,
	}

	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	if err != nil {
		return nil, nil, err
	}

	zkproof, err := groth16.Prove(cs, pk, witness)
	if err != nil {
		return nil, nil, err
	}

	publicWitness, err := witness.Public()
	if err != nil {
		return nil, nil, err
	}

	return zkproof, publicWitness, nil
}

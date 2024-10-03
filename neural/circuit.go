package circuit

import (
	"math/big"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/cmp"
)

type CircuitInputs struct {
	X []int `json:"input_vector"`
	W []int `json:"weights"`
	B int   `json:"bias"`
	Y int   `json:"output"`
}

type Circuit struct {
	// Inputs
	X [2]frontend.Variable `gnark:",public"` // Input vector (public for verification)
	W [2]frontend.Variable // Weights (private)
	B frontend.Variable    // Bias (private)

	// Output
	Y frontend.Variable `gnark:",public"` // Output after activation (public for verification)
}

// Define declares the circuit's constraints
func (circuit *Circuit) Define(api frontend.API) error {
	// Compute weighted sum: sum = w0*x0 + w1*x1 + b
	sum := api.Add(
		api.Mul(circuit.W[0], circuit.X[0]),
		api.Mul(circuit.W[1], circuit.X[1]),
		circuit.B,
	)

	// Implement ReLU activation: y = max(0, sum)
	// Since we can't have negative numbers in finite fields, we interpret numbers greater than (p-1)/2 as negative.

	// Get the modulus of the field
	modulus := api.Compiler().Field()
	halfModulus := new(big.Int).Rsh(modulus, 1) // (p - 1) / 2

	isPositive := cmp.IsLess(api, sum, halfModulus)

	y := api.Select(isPositive, sum, 0)

	api.AssertIsEqual(y, circuit.Y)

	return nil
}

func GenerateZKProof(cs constraint.ConstraintSystem, pk groth16.ProvingKey, customInputs any) (groth16.Proof, witness.Witness, error) {

	circuitInputs := customInputs.(CircuitInputs)

	// contruct assignment using dynamic veriables
	assignment := Circuit{
		X: [2]frontend.Variable{circuitInputs.X[0], circuitInputs.X[1]},
		W: [2]frontend.Variable{circuitInputs.W[0], circuitInputs.W[1]},
		B: circuitInputs.B,
		Y: circuitInputs.Y,
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

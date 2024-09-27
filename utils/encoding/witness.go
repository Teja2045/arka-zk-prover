package encoding

import (
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/witness"
)

func MarshalWitness(witness witness.Witness) ([]byte, error) {
	return witness.MarshalBinary()
}

func UnMarshalWitness(witnessBytes []byte) (witness.Witness, error) {
	// create new witness and unmarshal msg
	publicWitness, err := witness.New(ecc.BN254.ScalarField())
	if err != nil {
		return nil, err
	}

	err = publicWitness.UnmarshalBinary(witnessBytes)
	if err != nil {
		return nil, err
	}

	return publicWitness, nil
}

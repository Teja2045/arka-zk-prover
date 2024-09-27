package encoding

import (
	"encoding/json"

	"github.com/consensys/gnark/backend/groth16"
	bn254 "github.com/consensys/gnark/backend/groth16/bn254"
)

func MarhalZKProof(proof groth16.Proof) ([]byte, error) {
	proofbn254 := proof.(*bn254.Proof)
	proofBytes, err := json.Marshal(proofbn254)
	if err != nil {
		return nil, err
	}
	return proofBytes, nil
}

func UnMarshalZKProof(zkproof []byte) (groth16.Proof, error) {
	proof := new(bn254.Proof)
	err := json.Unmarshal(zkproof, proof)
	if err != nil {
		return nil, err
	}

	return proof, nil
}

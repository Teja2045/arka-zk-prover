package zkprover

import (
	"github.com/arka-labs/zk-prover/circuit"
	"github.com/arka-labs/zk-prover/utils/encoding"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
)

type ZKValidityProof struct {
	ZKProof       []byte `json:"zk_proof"`
	PublicWitness []byte `json:"public_witness"`
}

func NewZKValidityProof(zkProof groth16.Proof, publicWitness witness.Witness) (ZKValidityProof, error) {
	witnessBytes, err := encoding.MarshalWitness(publicWitness)
	if err != nil {
		return ZKValidityProof{}, err
	}

	zkProofBytes, err := encoding.MarhalZKProof(zkProof)
	if err != nil {
		return ZKValidityProof{}, err
	}
	return ZKValidityProof{
		ZKProof:       zkProofBytes,
		PublicWitness: witnessBytes,
	}, nil
}

func (validityProof *ZKValidityProof) Verify(dir string) error {
	zkProof, err := encoding.UnMarshalZKProof(validityProof.ZKProof)
	if err != nil {
		return err
	}

	publicWitness, err := encoding.UnMarshalWitness(validityProof.PublicWitness)
	if err != nil {
		return err
	}

	verfierKey, err := circuit.GetVerifierKey(dir)
	if err != nil {
		return err
	}

	return groth16.Verify(zkProof, verfierKey, publicWitness)
}

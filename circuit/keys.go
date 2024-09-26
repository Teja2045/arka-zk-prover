package circuit

import (
	"bytes"
	"fmt"

	"log/slog"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

func GenerateZKKeys() error {

	slog.Info(
		fmt.Sprintf(
			"Generating prover key, verifier key and constraint system for circuit",
		),
	)

	var circuit Circuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		return fmt.Errorf(fmt.Sprintln("error while compiling circuit: ", err))
	}
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		return fmt.Errorf(fmt.Sprintln("error while circuit setup", err))
	}

	vkBuf, pkBuf, ccsBuf := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)
	pk.WriteTo(pkBuf)
	vk.WriteTo(vkBuf)
	ccs.WriteTo(ccsBuf)

	proverKeyfileName := fmt.Sprintf("../keys/prover")
	WriteToFile(proverKeyfileName, pkBuf)

	verifierfileName := fmt.Sprintf("../keys/verifier")
	WriteToFile(verifierfileName, vkBuf)

	ccsfilename := fmt.Sprintf("../keys/ccs")
	WriteToFile(ccsfilename, ccsBuf)

	slog.Info("Keys are successfully generated and stored in {keys} folder")

	return nil

}

func WriteToFile(filename string, dataBuf *bytes.Buffer) {
	os.WriteFile(filename, dataBuf.Bytes(), 0666)
}

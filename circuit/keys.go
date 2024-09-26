package circuit

import (
	"bytes"
	"fmt"
	"path/filepath"

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
		return fmt.Errorf("error while compiling circuit: %v", err)
	}

	// Generate prover key and verifier key using groth16
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		return fmt.Errorf("error during circuit setup: %v", err)
	}

	// Prepare buffers for keys and constraint system
	vkBuf, pkBuf, ccsBuf := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)
	pk.WriteTo(pkBuf)
	vk.WriteTo(vkBuf)
	ccs.WriteTo(ccsBuf)

	// Define filenames
	proverKeyfileName := "./keys/prover"
	verifierfileName := "./keys/verifier"
	ccsfilename := "./keys/ccs"

	// Write keys and constraint system to files
	err = WriteToFile(proverKeyfileName, pkBuf)
	if err != nil {
		return err
	}
	err = WriteToFile(verifierfileName, vkBuf)
	if err != nil {
		return err
	}
	err = WriteToFile(ccsfilename, ccsBuf)
	if err != nil {
		return err
	}

	slog.Info("Keys are successfully generated and stored in the {keys} folder")
	return nil
}

// WriteToFile creates the directory if necessary and writes the buffer content to the specified file.
func WriteToFile(filename string, dataBuf *bytes.Buffer) error {
	// Extract directory path from the filename
	dir := filepath.Dir(filename)

	// Create directory if it doesn't exist
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %v", dir, err)
	}

	// Write data to file
	err = os.WriteFile(filename, dataBuf.Bytes(), 0666)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %v", filename, err)
	}

	return nil
}

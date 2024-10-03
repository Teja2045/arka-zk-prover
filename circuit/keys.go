package circuit

import (
	"bytes"
	"fmt"
	"path/filepath"

	"log/slog"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	bn254 "github.com/consensys/gnark/backend/groth16/bn254"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

// Define filenames
const (
	ProverKeyfileName = "prover"
	VerifierfileName  = "verifier"
	Csfilename        = "ccs"
)

func ProverKeyPath(dir string) string {
	return dir + "/" + ProverKeyfileName
}

func VerifierKeyPath(dir string) string {
	return dir + "/" + VerifierfileName
}

func CsPath(dir string) string {
	return dir + "/" + Csfilename
}

func GenerateZKKeys(dir string) error {

	proverKeyPath := ProverKeyPath(dir)
	verifierKeyPath := VerifierKeyPath(dir)
	csPath := CsPath(dir)

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

	// Write keys and constraint system to files
	err = WriteToFile(proverKeyPath, pkBuf)
	if err != nil {
		return err
	}
	err = WriteToFile(verifierKeyPath, vkBuf)
	if err != nil {
		return err
	}
	err = WriteToFile(csPath, ccsBuf)
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

// read the cs stored in the file and return it
func GetContraintSystem(dir string) (constraint.ConstraintSystem, error) {
	csPath := CsPath(dir)
	csBytes, err := os.ReadFile(csPath)
	if err != nil {
		return nil, err
	}

	cs := groth16.NewCS(ecc.BN254)
	_, err = cs.ReadFrom(bytes.NewBuffer(csBytes))
	if err != nil {
		return nil, err
	}

	return cs, nil
}

func GetProverKeyBytes(dir string) ([]byte, error) {
	proverKeyPath := ProverKeyPath(dir)
	return os.ReadFile(proverKeyPath)
}

// read the prover key stored in the file and return it
func GetProverKey(dir string) (groth16.ProvingKey, error) {
	proverKeyPath := ProverKeyPath(dir)
	pkBytes, err := os.ReadFile(proverKeyPath)
	if err != nil {
		return nil, err
	}

	pk := new(bn254.ProvingKey)
	_, err = pk.ReadFrom(bytes.NewBuffer(pkBytes))
	if err != nil {
		return nil, err
	}

	return pk, nil
}

func GetVerifierKeyBytes(dir string) ([]byte, error) {
	verifierKeyPath := VerifierKeyPath(dir)
	return os.ReadFile(verifierKeyPath)
}

// read the verifier key stored in the file and return it
func GetVerifierKey(dir string) (groth16.VerifyingKey, error) {
	verifierKeyPath := VerifierKeyPath(dir)
	vkBytes, err := os.ReadFile(verifierKeyPath)
	if err != nil {
		return nil, err
	}

	vk := new(bn254.VerifyingKey)
	_, err = vk.ReadFrom(bytes.NewBuffer(vkBytes))
	if err != nil {
		return nil, err
	}

	return vk, nil
}

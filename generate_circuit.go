package zkprover

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/arka-labs/zk-prover/circuit"
)

func GenerateKeys(dir string) error {
	return circuit.GenerateZKKeys(dir)
}

func GetZKProof(inputs any, dir string) (ZKValidityProof, error) {
	cs, err := circuit.GetContraintSystem(dir)
	if err != nil {
		return ZKValidityProof{}, err
	}

	pk, err := circuit.GetProverKey(dir)
	if err != nil {
		return ZKValidityProof{}, err
	}
	zkProof, publicWitness, err := circuit.GenerateZKProof(cs, pk, inputs)
	if err != nil {
		return ZKValidityProof{}, err
	}

	return NewZKValidityProof(zkProof, publicWitness)

}

// GenerateCircuit extracts the "circuit.go" file from a zip archive and saves it to "circuit/circuit.go", replacing any existing file.
func GenerateCircuit(r *zip.ReadCloser, destDir string) error {

	defer r.Close()

	// Iterate over the files in the zip archive
	for _, f := range r.File {
		// We are interested in the file named "circuit.go"
		if filepath.Base(f.Name) == "circuit.go" {
			// Open the file inside the zip archive
			rc, err := f.Open()
			if err != nil {
				return fmt.Errorf("failed to open file inside zip: %v", err)
			}
			defer rc.Close()

			// Create the destination directory if it doesn't exist
			destPath := filepath.Join(destDir, "circuit")
			err = os.MkdirAll(destPath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create destination directory: %v", err)
			}

			// Create (or overwrite) the "circuit.go" file in the "circuit/" directory
			destFile := filepath.Join(destPath, "circuit.go")
			outFile, err := os.Create(destFile)
			if err != nil {
				return fmt.Errorf("failed to create file: %v", err)
			}
			defer outFile.Close()

			// Copy the content of "circuit.go" from the zip to the destination file
			_, err = io.Copy(outFile, rc)
			if err != nil {
				return fmt.Errorf("failed to copy file content: %v", err)
			}

			fmt.Println("File circuit/circuit.go created/overwritten successfully")
			return nil
		}
	}
	return fmt.Errorf("circuit.go not found in the zip archive")
}

// GenerateCircuitFromRemoteURL downloads a zip file from the given URL, extracts the "circuit.go" file,
// and saves it to the "circuit/" directory, replacing any existing file.
func GenerateCircuitFromRemoteURL(zipFileURL, destDir string) error {
	// Create a temporary file to store the downloaded zip
	tmpZipFile, err := os.CreateTemp("", "circuit-*.zip")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpZipFile.Name()) // Clean up the temp file afterwards
	defer tmpZipFile.Close()

	// Download the zip file from the given URL
	resp, err := http.Get(zipFileURL)
	if err != nil {
		return fmt.Errorf("failed to download zip file from URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download zip file, got HTTP status: %v", resp.StatusCode)
	}

	// Copy the downloaded content to the temp zip file
	_, err = io.Copy(tmpZipFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save downloaded zip content: %v", err)
	}

	// Re-open the temp file for reading as a zip
	zipReader, err := zip.OpenReader(tmpZipFile.Name())
	if err != nil {
		return fmt.Errorf("failed to open downloaded zip file: %v", err)
	}
	defer zipReader.Close()

	// Call GenerateCircuit to extract and process the zip content
	return GenerateCircuit(zipReader, destDir)
}

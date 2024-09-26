package main

import (
	"archive/zip"
	"fmt"

	"log"

	zkprover "github.com/arka-labs/zk-prover"
	"github.com/arka-labs/zk-prover/circuit"
)

func main() {
	// Example usage
	zipFile := "example-circuit.zip" // Path to your zip file
	destDir := "../"
	// Open the zip file
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		log.Fatal("failed to open zip file:", err)
	}
	// Root directory (current directory)

	err = zkprover.GenerateCircuit(r, destDir)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("circuit.go successfully extracted and placed in the circuit directory")
	}

	circuit.GenerateZKKeys()
}

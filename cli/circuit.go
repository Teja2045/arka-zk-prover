package main

import (
	"archive/zip"
	"flag"
	"fmt"

	"log"

	zkprover "github.com/arka-labs/zk-prover"
)

func main() {
	// Define a command-line flag for the zip file path
	zipFile := flag.String("zipfile", "example-circuit.zip", "Path to your zip file")
	destDir := flag.String("dest", "./", "Destination directory for the circuit")
	flag.Parse()

	// Open the zip file
	r, err := zip.OpenReader(*zipFile)
	if err != nil {
		log.Fatal("Failed to open zip file:", err)
	}
	defer r.Close()

	// Call the GenerateCircuit function (assuming it's properly imported)
	err = zkprover.GenerateCircuit(r, *destDir)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("circuit.go successfully extracted and placed in the circuit directory")
	}
}

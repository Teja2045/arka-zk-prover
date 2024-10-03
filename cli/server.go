package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"log"

	zkprover "github.com/arka-labs/zk-prover"
	"github.com/arka-labs/zk-prover/server"
)

const (
	PORT  = 8000
	DIR_2 = "./keys"
)

func main() {

	// Define a command-line flag for the zip file path
	zipFile := flag.String("zipfile", "example-circuit.zip", "Path to your zip file")
	destDir := flag.String("dest", "./", "Destination directory for the circuit")
	flag.Parse()

	server := server.NewZKServer(8000, DIR_2)

	// Open the zip file
	r, err := zip.OpenReader(*zipFile)
	if err != nil {
		log.Fatal("failed to open zip file:", err)
	}
	// Root directory (current directory)

	err = zkprover.GenerateCircuit(r, *destDir)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("circuit.go successfully extracted and placed in the circuit directory")
	}

	err = zkprover.GenerateKeys(server.KeysDir)
	if err != nil {
		log.Fatal("Error while generating keys:", err)
	}

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

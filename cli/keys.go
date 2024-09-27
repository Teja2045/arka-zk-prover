package main

import "github.com/arka-labs/zk-prover/circuit"

const (
	DIR = "./keys"
)

func main() {
	circuit.GenerateZKKeys(DIR)
}

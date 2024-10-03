# Makefile

.PHONY: all run-circuit run-keys run-proof start-server

generate-example-circuit:
	go run cli/circuit.go

generate-neural-circuit:
	go run cli/circuit.go -zipfile="example-neural-circuit.zip"

generate-hash-circuit:
	go run cli/circuit.go -zipfile="hash-circuit.zip"

# Command to run cli/keys.go
generate-zk-keys:
	go run cli/keys.go

verify-neural-proof:
	go run cli/neural_proof.go

verify-hash-proof:
	go run cli/hash_proof.go

start-server:
	go run cli/server.go	

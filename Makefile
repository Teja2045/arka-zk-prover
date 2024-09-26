# Makefile

.PHONY: all run-circuit run-keys run-proof

generate-example-circuit:
	go run cli/circuit.go

# Command to run cli/keys.go
generate-zk-keys:
	go run cli/keys.go

verify-example-proof:
	go run cli/proof.go

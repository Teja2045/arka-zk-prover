package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	zkprover "github.com/arka-labs/zk-prover"
	"github.com/arka-labs/zk-prover/circuit"
)

func (server *ZKServer) Start() error {
	http.HandleFunc("/generate-zk-proof", server.handleZKProofRequest)

	addr := fmt.Sprintf(":%d", server.Port)
	log.Printf("Starting server on %s\n", addr)
	return http.ListenAndServe(addr, nil)
}

func (server *ZKServer) handleZKProofRequest(w http.ResponseWriter, r *http.Request) {
	// Parse input (you may need to customize this based on your input structure)
	var circuitInputs circuit.CircuitInputs
	err := json.NewDecoder(r.Body).Decode(&circuitInputs)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	// Call the zkprover to get the zkProof and publicWitness (mocked here)
	response, err := zkprover.GetZKProof(circuitInputs, server.KeysDir) // Assuming zkprover is an imported package
	if err != nil {
		http.Error(w, "failed to generate proof", http.StatusInternalServerError)
		return
	}

	// Set response header and send JSON response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

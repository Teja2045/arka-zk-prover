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
	http.HandleFunc("/verifier-key", server.handleGetVerifierKeyRequest)

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

func (server *ZKServer) handleGetVerifierKeyRequest(w http.ResponseWriter, _ *http.Request) {
	vkBytes, err := circuit.GetVerifierKeyBytes(server.KeysDir)
	if err != nil {
		http.Error(w, "failed to fetch verifier key", http.StatusInternalServerError)
		return
	}

	// Set the content type to application/octet-stream or any appropriate type.
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)

	// Write the verifier key bytes to the response
	_, writeErr := w.Write(vkBytes)
	if writeErr != nil {
		http.Error(w, "failed to write verifier key", http.StatusInternalServerError)
		return
	}
}

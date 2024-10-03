package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	zkprover "github.com/arka-labs/zk-prover"
	"github.com/arka-labs/zk-prover/circuit"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	bn254 "github.com/consensys/gnark/backend/groth16/bn254"
	"github.com/consensys/gnark/backend/witness"
)

func GetZKProof() (zkprover.ZKValidityProof, error) {
	// Define the endpoint URL
	url := "http://localhost:8000/generate-zk-proof" // Replace with your actual endpoint URL

	// Construct the CircuitInputs
	inputs := circuit.CircuitInputs{
		X: []int{3, 4},
		W: []int{2, -1},
		B: 1,
		Y: 3,
	}

	// Serialize the inputs to JSON
	jsonData, err := json.Marshal(inputs)
	if err != nil {
		return zkprover.ZKValidityProof{}, err
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return zkprover.ZKValidityProof{}, err
	}

	// Set the appropriate headers
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return zkprover.ZKValidityProof{}, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return zkprover.ZKValidityProof{}, err
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return zkprover.ZKValidityProof{}, err
	}

	// Deserialize the response into ZKValidityProof
	var proof zkprover.ZKValidityProof
	err = json.Unmarshal(body, &proof)
	if err != nil {
		return zkprover.ZKValidityProof{}, err
	}

	return proof, nil
}

func GetVerifierKey() ([]byte, error) {
	// Define the endpoint URL
	url := "http://localhost:8000/verifier-key" // Replace with your actual endpoint URL

	// Create a new HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("HTTP error: %s\nResponse body: %s\n", resp.Status, string(body))
		return nil, err
	}

	// Read the response body (verifier key bytes)
	return io.ReadAll(resp.Body)

}

func TestProof() {

	proof, err := GetZKProof()
	if err != nil {
		panic(err.Error())
	}

	zkproof, err := UnMarshalZKProof(proof.ZKProof)
	if err != nil {
		panic(err)
	}

	publicWitness, err := UnMarshalWitness(proof.PublicWitness)
	if err != nil {
		panic(err)
	}

	vk, err := GetVerifierKey()
	if err != nil {
		panic(err)
	}

	verifierKey, err := UnmarshalVerifierKey(vk)
	if err != nil {
		panic(err)
	}

	groth16.Verify(zkproof, verifierKey, publicWitness)

}

func UnMarshalZKProof(zkproof []byte) (groth16.Proof, error) {
	proof := new(bn254.Proof)
	err := json.Unmarshal(zkproof, proof)
	if err != nil {
		return nil, err
	}

	return proof, nil
}

func UnMarshalWitness(witnessBytes []byte) (witness.Witness, error) {
	// create new witness and unmarshal msg
	publicWitness, err := witness.New(ecc.BN254.ScalarField())
	if err != nil {
		return nil, err
	}

	err = publicWitness.UnmarshalBinary(witnessBytes)
	if err != nil {
		return nil, err
	}

	return publicWitness, nil
}

func UnmarshalVerifierKey(verifierKeyBytes []byte) (groth16.VerifyingKey, error) {

	verfierKey := new(bn254.VerifyingKey)
	_, err := verfierKey.ReadFrom(bytes.NewBuffer(verifierKeyBytes))
	if err != nil {
		return nil, err
	}

	return verfierKey, nil
}

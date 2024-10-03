# Arka-zk-prover

ZK prover for agents

### Imports
import the module
```go
import (
    //....
    zkprover "github.com/arka-labs/github.com/arka-labs/zk-prover"
    //...
)
```

### Generate circuit from .zip file
Call zkprover.GenerateCircuit method to generate circuit from .zip file 
```go
    func GenerateCircuit(r *zip.ReadCloser, destDir string) error {
        //....
    }
```

**Note:** this method only works if it is being called for the module itself. This package won't work via import.
To use different circuit, clone the project, replace `circuit/circuit.go` file with your own circuit logic. But make sure to follow the same format. You can refer to [this guide](./docs/01_circuit.md) on how to format the circuit logic.

### Generate constraint system, prover and verifier keys
call `zkprover.GenerateKeys()` to generate contraint system, prover and verfier keys and store them in `keys` folder

```go
func GenerateKeys() error {
    //....
}
```

### Get ZK Proof and public witness
call `zkprover.GetZKProof()` to generate zk proof for the given inputs.

```go
func GetZKProof(inputs ...any) (groth16.Proof, witness.Witness, error) {
    //....
}
```

## Circuits
Neural network circuit and Hashing circuit have been are defined `neural/circuit.go` and `gpt_hash/circuit.go files respectively`

### Run neural circuit
- run the following command to use neural circuit
```sh
make generate-neural-circuit
```
- run the following command to generate required keys for the circuit
```sh
make generate-zk-keys
```
- run the following command to generate zk proof and verify it for some example inputs
```sh
make verify-neural-proof
```

### Run sha256 hash circuit
- run the following command to use neural circuit
```sh
make generate-hash-circuit
```
- run the following command to generate required keys for the circuit
```sh
make generate-zk-keys
```
- run the following command to generate zk proof and verify it for some example inputs
```sh
make verify-hash-proof
```


## HTTP Server
Run http server to generate proofs

### Circuit
In `cli/server.go` file, replace the zipfile with your circuit.go zip file
```go
zipFile := "example-circuit.zip" // Path to your zip file
```

### Start Server
run the following command to start the server
```sh
make start-server
```

### APIs
the server exposes two apis

**Verifier key:** api to get the verifier key which can be used to verify the zk proof

type: `GET` request
endpoint: `/verifier-key`
request body: nil
response: bytes data of verifier key

**generate ZK Proof:** api to generate the zkproof which included both zk proof and public witness

type: GET request
endpoint: /generate-zk-proof
request body: circuit inputs json
response: {
    "zk_proof": bytes data,
    "public_witness": bytes data
}


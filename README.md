# arka-zk-prover

ZK prover for agents

### Generate circuit from .zip file
Call zkProver.GenerateCircuit method to generate circuit from .zip file 
```go
    func GenerateCircuit(r *zip.ReadCloser, destDir string) error {
        //....
    }
```

### Generate constraint system, prover and verifier keys
call zkProver.GenerateKeys() to generate contraint system, prover and verfier keys and store them in files

```go
func GenerateKeys() error {
    //....
}
```

### Get ZK Proof and public witness
call zkProver.GetZKProof() to get zk proof out of given inputs.

```go
func GetZKProof(inputs ...any) (groth16.Proof, witness.Witness, error) {
    //....
}
```

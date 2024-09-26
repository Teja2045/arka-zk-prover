# arka-zk-prover

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

### Generate constraint system, prover and verifier keys
call zkprover.GenerateKeys() to generate contraint system, prover and verfier keys and store them in files

```go
func GenerateKeys() error {
    //....
}
```

### Get ZK Proof and public witness
call zkprover.GetZKProof() to get zk proof out of given inputs.

```go
func GetZKProof(inputs ...any) (groth16.Proof, witness.Witness, error) {
    //....
}
```

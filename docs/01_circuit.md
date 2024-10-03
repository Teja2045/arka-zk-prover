# Guide to write the circuit
Follow the below format of defining circuit to maker sure the other apis defined in the module packages can be used seamlessly.

**Table of Contents**

1. [Introduction](#introduction)
2. [Defining the Circuit Struct](#defining-the-circuit-struct)
3. [Defining the CircuitInputs Struct](#defining-the-circuitinputs-struct)
4. [Implementing the Define Method](#implementing-the-define-method)
5. [Generating the ZK Proof](#generating-the-zk-proof)
6. [Conclusion](#conclusion)



## Introduction

Gnark is a Go library for writing zero-knowledge proofs (ZKPs) using zk-SNARKs. When creating a circuit in Gnark, it's essential to structure your code correctly to define the variables, constraints, and proof generation logic. This guide provides a step-by-step approach to formatting a circuit in Gnark, covering the key components:

- Defining the circuit struct
- Defining the circuit inputs struct
- Implementing the `Define` method
- Writing the `GenerateZKProof` function

## Defining the Circuit Struct

The circuit struct represents the variables and constraints of your zero-knowledge proof. Each field corresponds to an input, output, or intermediate variable in the circuit.

```go
type Circuit struct {
   // needed inputs
}
```

**Example**:
```go
type Circuit struct {
    // Public Inputs
    A frontend.Variable `gnark:",public"` // Public input A
    B frontend.Variable `gnark:",public"` // Public input B

    // Private Inputs
    C frontend.Variable // Private input C
    D frontend.Variable // Private input D

}

```

**Key Points:**

- **Public Variables:** Annotate public variables with the struct tag `` `gnark:",public"` ``.
- **Variable Types:** Use `frontend.Variable` for all variables.
- **Naming Convention:** Make sure name the struct Circuit.

## Defining the CircuitInputs Struct

The `CircuitInputs` struct is used to provide concrete values to the circuit variables when generating a proof or verifying one.

```go
type CircuitInputs struct {
   field1 datatype `json:"json_field_name"`
}
```

**Key Points:**

- **Field Matching:** Ensure the fields match the variables in the `Circuit` struct.
- **JSON Tags:** Use JSON tags if you plan to serialize/deserialize the inputs(used for http endpoints).
- **Data Types:** Choose appropriate data types (e.g., `int`).
- **Struct Name:** Make sure to the struct `CircuitInputs`.

## Implementing the Define Method

The `Define` method specifies the constraints of the circuit. This is where you define the mathematical relationships between the variables.

```go
func (circuit *Circuit) Define(api frontend.API) error {
    // circuit logic
}
```

**Example:**
```go
func (circuit *Circuit) Define(api frontend.API) error {
    // Example Constraint: C = (A + B)

    // Compute sum = A + B
    sum := api.Add(circuit.A, circuit.B)
   // assert 
    api.AssertIsEqual(circuit.C, sum)

    return nil
}

```

**Key Points:**

- **Arithmetic Operations:** Use `api.Add`, `api.Sub`, `api.Mul`, etc., for operations.
- **Constraints:** Use `api.AssertIsEqual` to enforce constraints.
- **Error Handling:** Return `nil` if no errors occur.

## Generating the ZK Proof

The `GenerateZKProof` function generates a zero-knowledge proof using the circuit, proving key, and inputs. It is used by other method in the module. So make sure to use the same function name and signature.

**Example:**

```go
func GenerateZKProof(
    cs constraint.ConstraintSystem,
    pk groth16.ProvingKey,
    customInputs any,
) (groth16.Proof, witness.Witness, error) {
    // Type assertion to convert customInputs to CircuitInputs
    circuitInputs := customInputs.(CircuitInputs)

    // Create an assignment of variables with provided inputs
    assignment := Circuit{
        A: circuitInputs.A,
        B: circuitInputs.B,
        C: circuitInputs.C,
    }

    // Create a witness for the circuit
    witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
    if err != nil {
        return nil, nil, err
    }

    // Generate the proof
    proof, err := groth16.Prove(cs, pk, witness)
    if err != nil {
        return nil, nil, err
    }

    // Extract the public part of the witness
    publicWitness, err := witness.Public()
    if err != nil {
        return nil, nil, err
    }

    return proof, publicWitness, nil
}
```



**Key Points:**
- **assignment:** the assignment is the only part that needs to changed. Everything else can be reused for all circuits (there could be rare exceptions which requires more customization).

## Conclusion

Formatting a circuit in Gnark involves defining the circuit and inputs structs, implementing the `Define` method to specify constraints, and writing a function to generate the zero-knowledge proof. By following this structure, you ensure that your circuit is correctly set up to be used by other apis in the module packages.

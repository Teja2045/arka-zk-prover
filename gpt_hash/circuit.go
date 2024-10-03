package circuit

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash"
	"github.com/consensys/gnark/std/math/bits"
	"github.com/consensys/gnark/std/math/bitslice"
	"github.com/consensys/gnark/std/math/cmp"
	"github.com/consensys/gnark/std/math/uints"
	"github.com/consensys/gnark/std/permutation/sha2"
)

// var _ hash.BinaryFixedLengthHasher = (*digest)(nil)
// var _ hash.FieldHasher = (*digest)(nil)

type CircuitInputs struct {
	Data []byte   `json:"data"`
	Hash [32]byte `json:"hash"`
}

type Circuit struct {
	In       []uints.U8
	Expected []uints.U8
}

func (circuit *Circuit) Define(api frontend.API) error {
	hFunc, err := New(api)
	if err != nil {
		return err
	}

	uapi, err := uints.New[uints.U32](api)
	if err != nil {
		return err
	}
	hFunc.Write(circuit.In)
	hash := hFunc.Sum()
	if len(hash) != 32 {
		return fmt.Errorf("not 32 bytes")
	}

	for i := range circuit.Expected {
		uapi.ByteAssertEq(circuit.Expected[i], hash[i])
	}
	return nil
}

func GenerateZKProof(cs constraint.ConstraintSystem, pk groth16.ProvingKey, customInputs any) (groth16.Proof, witness.Witness, error) {

	circuitInputs := customInputs.(CircuitInputs)

	// contruct assignment using dynamic veriables
	assignment := Circuit{
		In:       uints.NewU8Array(circuitInputs.Data),
		Expected: uints.NewU8Array(circuitInputs.Hash[:]),
	}

	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	if err != nil {
		return nil, nil, err
	}

	zkproof, err := groth16.Prove(cs, pk, witness)
	if err != nil {
		return nil, nil, err
	}

	publicWitness, err := witness.Public()
	if err != nil {
		return nil, nil, err
	}

	return zkproof, publicWitness, nil
}

// -------------------------- Sha256 Implementation for Circuit ----------------------------

var _seed = uints.NewU32Array([]uint32{
	0x6A09E667, 0xBB67AE85, 0x3C6EF372, 0xA54FF53A, 0x510E527F, 0x9B05688C, 0x1F83D9AB, 0x5BE0CD19,
})

type digest struct {
	api  frontend.API
	uapi *uints.BinaryField[uints.U32]
	in   []uints.U8
}

func New(api frontend.API) (hash.BinaryFixedLengthHasher, error) {
	uapi, err := uints.New[uints.U32](api)
	if err != nil {
		return nil, err
	}
	return &digest{api: api, uapi: uapi}, nil
}

func (d *digest) Write(data []uints.U8) {
	d.in = append(d.in, data...)
}

func (d *digest) padded(bytesLen int) []uints.U8 {
	zeroPadLen := 55 - bytesLen%64
	if zeroPadLen < 0 {
		zeroPadLen += 64
	}
	if cap(d.in) < len(d.in)+9+zeroPadLen {
		// in case this is the first time this method is called increase the
		// capacity of the slice to fit the padding.
		d.in = append(d.in, make([]uints.U8, 9+zeroPadLen)...)
		d.in = d.in[:len(d.in)-9-zeroPadLen]
	}
	buf := d.in
	buf = append(buf, uints.NewU8(0x80))
	buf = append(buf, uints.NewU8Array(make([]uint8, zeroPadLen))...)
	lenbuf := make([]uint8, 8)
	binary.BigEndian.PutUint64(lenbuf, uint64(8*bytesLen))
	buf = append(buf, uints.NewU8Array(lenbuf)...)
	return buf
}

func (d *digest) Sum() []uints.U8 {
	var runningDigest [8]uints.U32
	var buf [64]uints.U8
	copy(runningDigest[:], _seed)
	padded := d.padded(len(d.in))
	for i := 0; i < len(padded)/64; i++ {
		copy(buf[:], padded[i*64:(i+1)*64])
		runningDigest = sha2.Permute(d.uapi, runningDigest, buf)
	}
	var ret []uints.U8
	for i := range runningDigest {
		ret = append(ret, d.uapi.UnpackMSB(runningDigest[i])...)
	}
	return ret
}

func (d *digest) FixedLengthSum(length frontend.Variable) []uints.U8 {
	// we need to do two things here -- first the padding has to be put to the
	// right place. For that we need to know how many blocks we have used. We
	// need to fit at least 9 more bytes (padding byte and 8 bytes for input
	// length). Knowing the block, we have to keep running track if the current
	// block is the expected one.
	//
	// idea - have a mask for blocks where 1 is only for the block we want to
	// use.

	data := make([]uints.U8, len(d.in))
	copy(data, d.in)

	comparator := cmp.NewBoundedComparator(d.api, big.NewInt(int64(len(data)+64+8)), false)

	for i := 0; i < 64+8; i++ {
		data = append(data, uints.NewU8(0))
	}

	lenMod64 := d.mod64(length)
	lenMod64Less56 := comparator.IsLess(lenMod64, 56)

	paddingCount := d.api.Sub(64, lenMod64)
	paddingCount = d.api.Select(lenMod64Less56, paddingCount, d.api.Add(paddingCount, 64))

	totalLen := d.api.Add(length, paddingCount)
	last8BytesPos := d.api.Sub(totalLen, 8)

	var dataLenBtyes [8]frontend.Variable
	d.bigEndianPutUint64(dataLenBtyes[:], d.api.Mul(length, 8))

	for i := range data {
		isPaddingStartPos := d.api.IsZero(d.api.Sub(i, length))
		data[i].Val = d.api.Select(isPaddingStartPos, 0x80, data[i].Val)

		isPaddingPos := comparator.IsLess(length, i)
		data[i].Val = d.api.Select(isPaddingPos, 0, data[i].Val)
	}

	for i := range data {
		isLast8BytesPos := d.api.IsZero(d.api.Sub(i, last8BytesPos))
		for j := 0; j < 8; j++ {
			if i+j < len(data) {
				data[i+j].Val = d.api.Select(isLast8BytesPos, dataLenBtyes[j], data[i+j].Val)
			}
		}
	}

	var runningDigest [8]uints.U32
	var resultDigest [8]uints.U32
	var buf [64]uints.U8
	copy(runningDigest[:], _seed)
	copy(resultDigest[:], _seed)

	for i := 0; i < len(data)/64; i++ {
		copy(buf[:], data[i*64:(i+1)*64])
		runningDigest = sha2.Permute(d.uapi, runningDigest, buf)

		isInRange := comparator.IsLess(i*64, totalLen)

		for j := 0; j < 8; j++ {
			for k := 0; k < 4; k++ {
				resultDigest[j][k].Val = d.api.Select(isInRange, runningDigest[j][k].Val, resultDigest[j][k].Val)
			}
		}
	}

	var ret []uints.U8
	for i := range resultDigest {
		ret = append(ret, d.uapi.UnpackMSB(resultDigest[i])...)
	}
	return ret
}

func (d *digest) Reset() {
	d.in = nil
}

func (d *digest) Size() int { return 32 }

func (d *digest) mod64(v frontend.Variable) frontend.Variable {
	lower, _ := bitslice.Partition(d.api, v, 6, bitslice.WithNbDigits(64))
	return lower
}

func (d *digest) bigEndianPutUint64(b []frontend.Variable, x frontend.Variable) {
	bts := bits.ToBinary(d.api, x, bits.WithNbDigits(64))
	for i := 0; i < 8; i++ {
		b[i] = bits.FromBinary(d.api, bts[(8-i-1)*8:(8-i)*8])
	}
}

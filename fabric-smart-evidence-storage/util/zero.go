package util

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	mimc2 "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/hash/mimc"
	
)

var pk groth16.ProvingKey

func InitGroth16(pkBytes []byte) {
	var buf bytes.Buffer
	pk = groth16.NewProvingKey(ecc.BN254)
	buf.Write(pkBytes)
	_, err := pk.ReadFrom(&buf)
	if err != nil {
		panic(fmt.Errorf("failed to read pk: %v", err))
	}

}

type UserCredential struct {
	Role frontend.Variable `gnark:"role"`
	Hash frontend.Variable `gnark:",public"`
}

func (circuit *UserCredential) Define(api frontend.API) error {
	mc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	api.AssertIsEqual(circuit.Role, "2")
	mc.Write(circuit.Role)
	api.AssertIsEqual(circuit.Hash, mc.Sum())
	return nil
}

func GenUserCredential(role *big.Int) ([]byte, []byte, error) {
	// 零知识证明 生成凭证
	var circuit UserCredential
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		return nil, nil, err
	}
	assignment := UserCredential{
		Role: role,
	}
	assignment.Hash, err = GenHash(role)
	if err != nil {
		return nil, nil, err
	}
	w1, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	if err != nil {
		return nil, nil, err
	}
	proof, err := groth16.Prove(ccs, pk, w1)
	if err != nil {
		return nil, nil, err
	}
	publicWitness, err := w1.Public()
	if err != nil {
		return nil, nil, err
	}
	var publicWitnessBytes bytes.Buffer
	_, err = publicWitness.WriteTo(&publicWitnessBytes)
	if err != nil {
		return nil, nil, err
	}
	var proofBytes bytes.Buffer
	_, err = proof.WriteTo(&proofBytes)
	if err != nil {
		return nil, nil, err
	}
	return publicWitnessBytes.Bytes(), proofBytes.Bytes(), nil
}

func GenHash(role *big.Int) ([]byte, error) {
	h := mimc2.NewMiMC()
	var b1 fr.Element
	b1.SetBytes(role.Bytes())
	_, err := h.Write(b1.Marshal())
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

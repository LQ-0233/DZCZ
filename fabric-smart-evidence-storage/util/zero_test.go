package util

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"log"
	"math/big"
	"os"
	"testing"
)

func write(name string, p []byte) {

	filePath := "./" + name
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}

	defer file.Close()
	write := bufio.NewWriter(file)
	write.Write(p)
	write.Flush()
}

func TestGen(t *testing.T) {
	var user UserCredential
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &user)
	if err != nil {
		log.Fatalf("Failed to compile circuit: %v", err)
	}
	// groth16 zkSNARK: Setup
	pk, vk, err := groth16.Setup(ccs)

	if err != nil {
		log.Fatalf("Failed to generate: %v", err)
	}

	var buf bytes.Buffer
	_, err = pk.WriteTo(&buf)
	if err != nil {
		log.Fatalf("Failed to write pk: %v", err)
	}
	write("./pk", buf.Bytes())

	buf.Reset()
	_, err = vk.WriteTo(&buf)
	if err != nil {
		log.Fatalf("Failed to write pk: %v", err)
	}
	write("./vk", buf.Bytes())
}

func TestGenUserCredential(t *testing.T) {
	pkBytes, err := os.ReadFile("../pk")
	if err != nil {
		panic(err)
	}
	InitGroth16(pkBytes)
	role := new(big.Int).SetInt64(2)
	userCredential, proof, err := GenUserCredential(role)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(base64.StdEncoding.EncodeToString(userCredential))
	fmt.Println(base64.StdEncoding.EncodeToString(proof))
}

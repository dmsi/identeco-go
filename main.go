package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math/big"
)

func decodeE(s string) (uint32, error) {
	bytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(append(bytes, 0)), nil
}

func decodeN(s string) (*big.Int, error) {
	// bytes, err := base64.URLEncoding.DecodeString(s)
	bytes, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}

	n := new(big.Int)
	n.SetBytes(bytes)
	return n, nil
}

func main() {
	e, _ := decodeE("AQAB")
	n, err := decodeN("tRXzVqY51HMCh-iK2K0YmGF044P2qM_42MDBZuk6CpqUg1Vm7ylBHLm41QWNIwvzyVtBiibjSPtT_Ua2-_6v5dz2bwZqUzxYU_yq5sacv3yfOpwe8mYej2wyaC0fBcKSigrpFj3nDHTXEUGIiR0Vptd7ja7vjOcj_8raGjaR7zGF_5P42OA-UUDmRmyU1PG_d4fV-bagip1byEcPM4GSxqOnWkJdNX9da82S9QxYSofFq9t8MYH2texM5ImcqZ0FmdUXb8k1DeBXv0dqg1ZbhaDvCzNWfgoMjhPeB5lpnCP0gR-X_3dLJDPI1lU0ddnjepCWuh48WuImxfilaoQCcw")
	fmt.Printf("e: %v, n: %v, err: %v\n", e, n, err)

	keypair, err := rsa.GenerateKey(rand.Reader, 256)
	s := base64.RawURLEncoding.EncodeToString(keypair.N.Bytes())
	fmt.Printf("s: %v\n", s)
}

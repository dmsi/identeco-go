package util

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
)

func KeyID(publicKey *rsa.PublicKey) string {
	hash := sha256.Sum256(publicKey.N.Bytes())
	kid := hex.EncodeToString(hash[:])
	return kid
}

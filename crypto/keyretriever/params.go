//go:build darwin || linux

package keyretriever

import (
	"hash"

	"src/crypto"
)

// pbkdf2Params holds platform-specific PBKDF2 key derivation parameters.
// Each platform file defines its own params variable.
type pbkdf2Params struct {
	salt       []byte
	iterations int
	keySize    int
	hashFunc   func() hash.Hash
}

// deriveKey derives an encryption key from a secret using PBKDF2.
func (p pbkdf2Params) deriveKey(secret []byte) []byte {
	return crypto.PBKDF2Key(secret, p.salt, p.iterations, p.keySize, p.hashFunc)
}

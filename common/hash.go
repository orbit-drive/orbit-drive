package common

import (
	"crypto/sha256"
)

// String parsing helpers
func HashStr(p string) []byte {
	hash := sha256.Sum256(ToByte(p))
	return hash[:]
}

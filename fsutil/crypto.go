package fsutil

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// HashStr parsing helpers
func HashStr(p string) []byte {
	return HashBytes(ToByte(p))
}

func HashBytes(b []byte) []byte {
	hash := sha256.Sum256(b)
	return hash[:]
}

// SecureHash
func SecureHash(p string) ([]byte, error) {
	return bcrypt.GenerateFromPassword(ToByte(p), bcrypt.DefaultCost)
}

// Md5Checksum calculates the md5 checksum of a file
// Todo: Read file in stream so file does not load entirely in memory.
func Md5Checksum(p string) (string, error) {
	file, err := os.Open(p)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := md5.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", nil
	}
	sum := hasher.Sum(nil)[:16]
	return hex.EncodeToString(sum), nil
}

package crypto

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
)

// HashType represents the type of hash algorithm
type HashType int

const (
	SHA256 HashType = iota
	SHA512
)

// HashString hashes a string using the specified algorithm
func HashString(data string, hashType HashType) (string, error) {
	var h hash.Hash
	
	switch hashType {
	case SHA256:
		h = sha256.New()
	case SHA512:
		h = sha512.New()
	default:
		return "", fmt.Errorf("unsupported hash type: %v", hashType)
	}
	
	_, err := h.Write([]byte(data))
	if err != nil {
		return "", err
	}
	
	return hex.EncodeToString(h.Sum(nil)), nil
}

// HashBytes hashes byte data using the specified algorithm
func HashBytes(data []byte, hashType HashType) (string, error) {
	var h hash.Hash
	
	switch hashType {
	case SHA256:
		h = sha256.New()
	case SHA512:
		h = sha512.New()
	default:
		return "", fmt.Errorf("unsupported hash type: %v", hashType)
	}
	
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}
	
	return hex.EncodeToString(h.Sum(nil)), nil
}

// HashWithSalt hashes data with a salt
func HashWithSalt(data, salt string, hashType HashType) (string, error) {
	return HashString(data+salt, hashType)
}

// VerifyHash verifies if data matches the hash
func VerifyHash(data, expectedHash string, hashType HashType) (bool, error) {
	actualHash, err := HashString(data, hashType)
	if err != nil {
		return false, err
	}
	return actualHash == expectedHash, nil
}

// VerifyHashWithSalt verifies if data with salt matches the hash
func VerifyHashWithSalt(data, salt, expectedHash string, hashType HashType) (bool, error) {
	actualHash, err := HashWithSalt(data, salt, hashType)
	if err != nil {
		return false, err
	}
	return actualHash == expectedHash, nil
}

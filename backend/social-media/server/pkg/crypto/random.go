package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"math/big"
)

// GenerateRandomBytes generates cryptographically secure random bytes
func GenerateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	return bytes, err
}

// GenerateRandomString generates a cryptographically secure random string
func GenerateRandomString(length int) (string, error) {
	bytes, err := GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// GenerateRandomHex generates a random hex string
func GenerateRandomHex(length int) (string, error) {
	bytes, err := GenerateRandomBytes(length / 2)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateRandomInt generates a random integer between min and max (inclusive)
func GenerateRandomInt(min, max int) (int, error) {
	if min >= max {
		return 0, errors.New("min must be less than max")
	}
	
	bigMax := big.NewInt(int64(max - min + 1))
	n, err := rand.Int(rand.Reader, bigMax)
	if err != nil {
		return 0, err
	}
	
	return int(n.Int64()) + min, nil
}

// GenerateSecureToken generates a secure token for various purposes
func GenerateSecureToken(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	token := make([]byte, length)
	
	for i := range token {
		num, err := GenerateRandomInt(0, len(charset)-1)
		if err != nil {
			return "", err
		}
		token[i] = charset[num]
	}
	
	return string(token), nil
}

// GenerateAPIKey generates a secure API key
func GenerateAPIKey() (string, error) {
	bytes, err := GenerateRandomBytes(32)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// GenerateSessionID generates a secure session ID
func GenerateSessionID() (string, error) {
	return GenerateRandomHex(32)
}

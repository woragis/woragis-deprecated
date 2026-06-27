package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	password := "TestPassword123!"
	hashed, err := HashPassword(password, 10)

	require.NoError(t, err)
	assert.NotEmpty(t, hashed)
	assert.NotEqual(t, password, hashed)
}

func TestVerifyPassword_Success(t *testing.T) {
	password := "TestPassword123!"
	hashed, err := HashPassword(password, 10)
	require.NoError(t, err)

	err = VerifyPassword(password, hashed)
	assert.NoError(t, err)
}

func TestVerifyPassword_Failure(t *testing.T) {
	password := "TestPassword123!"
	wrongPassword := "WrongPassword123!"
	hashed, err := HashPassword(password, 10)
	require.NoError(t, err)

	err = VerifyPassword(wrongPassword, hashed)
	assert.Error(t, err)
}

func TestCheckPasswordStrength_Success(t *testing.T) {
	validPasswords := []string{
		"SecurePass123!",
		"P@ssw0rd!",
		"MyStr0ng#Pass",
		"Test1234$",
	}

	for _, password := range validPasswords {
		err := CheckPasswordStrength(password)
		assert.NoError(t, err, "Password %s should be valid", password)
	}
}

func TestCheckPasswordStrength_TooShort(t *testing.T) {
	err := CheckPasswordStrength("Short1!")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least 8 characters")
}

func TestCheckPasswordStrength_NoUppercase(t *testing.T) {
	err := CheckPasswordStrength("password123!")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "uppercase")
}

func TestCheckPasswordStrength_NoLowercase(t *testing.T) {
	err := CheckPasswordStrength("PASSWORD123!")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "lowercase")
}

func TestCheckPasswordStrength_NoNumber(t *testing.T) {
	err := CheckPasswordStrength("Password!")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "number")
}

func TestCheckPasswordStrength_NoSpecialChar(t *testing.T) {
	err := CheckPasswordStrength("Password123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "special character")
}

func BenchmarkHashPassword(b *testing.B) {
	password := "TestPassword123!"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = HashPassword(password, 10)
	}
}

func BenchmarkVerifyPassword(b *testing.B) {
	password := "TestPassword123!"
	hashed, _ := HashPassword(password, 10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = VerifyPassword(password, hashed)
	}
}


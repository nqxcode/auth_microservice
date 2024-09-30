package utils

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// GenerateSalt generate salt for password hashing
func GenerateSalt(_ context.Context) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}

// HashPassword calculate hash for password with salt
func HashPassword(_ context.Context, password, salt string) (string, error) {
	saltedPassword := password + salt
	hash := sha256.Sum256([]byte(saltedPassword))
	hashedPassword, err := bcrypt.GenerateFromPassword(hash[:], bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword check password with hash
func VerifyPassword(_ context.Context, password, salt, hash string) bool {
	saltedPassword := password + salt
	hashToCheck := sha256.Sum256([]byte(saltedPassword))
	err := bcrypt.CompareHashAndPassword([]byte(hash), hashToCheck[:])
	return err == nil
}

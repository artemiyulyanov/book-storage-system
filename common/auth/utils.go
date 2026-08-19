package auth

import (
	"crypto/subtle"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func PasswordsHashEqual(password_hashed, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password_hashed), []byte(password))
	return err == nil
}

func PasswordsMatch(password, confirmPassword string) bool {
	return subtle.ConstantTimeCompare([]byte(password), []byte(confirmPassword)) == 1
}

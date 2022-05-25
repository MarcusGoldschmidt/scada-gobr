package auth

import (
	"bytes"
	"crypto/sha512"
)

func ValidatePassword(password string, passwordHash string) bool {
	check := sha512.Sum512([]byte(password))

	return bytes.Compare(check[:], []byte(passwordHash)) == 1
}

func MakeHash(password string) string {
	check := sha512.Sum512([]byte(password))

	return string(check[:])
}

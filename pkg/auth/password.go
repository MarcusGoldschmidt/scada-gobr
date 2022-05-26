package auth

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
)

func ValidatePassword(password string, passwordHash string) (bool, error) {
	check := sha512.Sum512([]byte(password))

	passwordHashBytes, err := base64.StdEncoding.DecodeString(passwordHash)
	if err != nil {
		return false, err
	}

	return bytes.Compare(check[:], passwordHashBytes) == 0, nil
}

func MakeHash(password string) string {
	check := sha512.Sum512([]byte(password))

	return base64.StdEncoding.EncodeToString(check[:])
}

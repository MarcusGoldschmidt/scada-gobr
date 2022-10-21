package auth

import "testing"

func TestHash(t *testing.T) {
	password := "teste"

	hash, err := MakeHash(password)

	if err != nil {
		t.Error(err)
	}

	validatePassword := ValidatePassword(password, hash)

	if !validatePassword {
		t.Fail()
	}
}

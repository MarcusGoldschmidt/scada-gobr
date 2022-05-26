package auth

import "testing"

func TestHash(t *testing.T) {
	password := "teste"

	hash := MakeHash(password)

	validatePassword, err := ValidatePassword(password, hash)
	if err != nil {
		t.Error(err)
	}

	if !validatePassword {
		t.Fail()
	}
}

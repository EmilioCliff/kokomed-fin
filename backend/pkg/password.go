package pkg

import (
	"golang.org/x/crypto/bcrypt"
)

func GenerateHashPassword(password string, cost int) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", Errorf(INTERNAL_ERROR, "failed to generate hash password: %s", err.Error())
	}

	return string(hashPassword), nil
}

func ComparePasswordAndHash(hashPass string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(password))
	if err != nil {
		return Errorf(AUTHENTICATION_ERROR, "invalid password")
	}

	return nil
}

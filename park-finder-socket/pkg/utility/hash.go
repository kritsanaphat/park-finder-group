package utility

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password []byte) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(hashedPassword []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

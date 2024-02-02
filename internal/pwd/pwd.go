package pwd

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gophermart/internal/entity"
)

type Pwd struct{}

func New() *Pwd {
	return &Pwd{}
}

func (ah *Pwd) Hash(pwd entity.UserPassword) (entity.UserPasswordHash, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Pwd - Pwd - bcrypt.GenerateFromPassword: %w", err)
	}
	return entity.UserPasswordHash(hash), nil
}

func (ah *Pwd) IsCorrect(pwd entity.UserPassword, hash entity.UserPasswordHash) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err == nil
}

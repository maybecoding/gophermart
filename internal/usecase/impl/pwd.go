package impl

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gophermart/internal/entity"
)

type PwdImpl struct {
}

func NewPwd() *PwdImpl {
	return &PwdImpl{}
}

func (ah *PwdImpl) Hash(pwd entity.UserPassword) (entity.UserPasswordHash, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("PwdImpl - PwdImpl - bcrypt.GenerateFromPassword: %w", err)
	}
	return entity.UserPasswordHash(hash), nil
}

func (ah *PwdImpl) IsCorrect(pwd entity.UserPassword, hash entity.UserPasswordHash) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err == nil
}

package usecase

import (
	"context"
	"errors"
	"fmt"
	"gophermart/internal/entity"
)

type AuthUseCase struct {
	repo AuthRepo
	pwd  Pwd
	jwt  Jwt
}

func NewAuth(r AuthRepo, pwd Pwd, jwt Jwt) *AuthUseCase {
	return &AuthUseCase{repo: r, pwd: pwd, jwt: jwt}
}

// Register - регистрация пользователя
//  1. Проверяем что имя пользователя доступно для добавления
//  2. Хэшируем пароль
//  3. Сохраняем пользователя
//  4. Формируем jwt-токен
func (uc *AuthUseCase) Register(ctx context.Context, login entity.UserLogin, pwd entity.UserPassword) (entity.Token, error) {
	// 1. Проверяем что имя пользователя доступно для добавления
	loginAvail, err := uc.repo.LoginAvailable(ctx, login)
	if err != nil {
		return "", fmt.Errorf("AuthUseCase - Register - repo.LoginAvailable: %w", err)
	}
	if !loginAvail {
		return "", entity.ErrUserNotAvailable
	}

	// 2. Хэшируем пароль
	pwdHash, err := uc.pwd.Hash(pwd)
	if err != nil {
		return "", fmt.Errorf("AuthUseCase - Register - pwd.Hash: %w", err)
	}

	// 3. Сохраняем пользователя
	user, err := uc.repo.UserNew(ctx, login, pwdHash)
	if err != nil {
		return "", fmt.Errorf("AuthUseCase - Register - repo.UserNew: %w", err)
	}
	if user == nil {
		return "", errors.New("AuthUseCase - Register - user isn't created")
	}

	// 4. Формируем tokenData-токен
	tokenData := entity.TokenData{UserID: user.UserID}
	jwtToken, err := uc.jwt.Get(tokenData)
	if err != nil {
		return "", fmt.Errorf("AuthUseCase - Register - pwd.Get: %w", err)
	}
	return jwtToken, nil
}

// Login - Аутентификация пользователя
// 1. Получаем по логину UserID и Hash из БД
// 2. Проверяем Хэш
// 3. Формируем jwt-токен
func (uc *AuthUseCase) Login(ctx context.Context, login entity.UserLogin, pwd entity.UserPassword) (entity.Token, error) {
	// 1. Получаем по логину UserID и Hash из БД
	user, err := uc.repo.UserGet(ctx, login)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			return "", entity.ErrUserNotFound
		}
		return "", fmt.Errorf("AuthUseCase - Login - repo.UserFind: %w", err)
	}
	if user == nil {
		return "", entity.ErrUserNotFound
	}

	// 2. Проверяем Хэш
	pwdOk := uc.pwd.IsCorrect(pwd, user.UserPasswordHash)
	if !pwdOk {
		return "", entity.ErrIncorrectPassword
	}

	// 3. Формируем jwt-токен
	tokenData := entity.TokenData{UserID: user.UserID}
	jwtToken, err := uc.jwt.Get(tokenData)
	if err != nil {
		return "", fmt.Errorf("AuthUseCase - Register - pwd.Get: %w", err)
	}
	return jwtToken, nil
}

func (uc *AuthUseCase) ValidateToken(t entity.Token) (entity.TokenData, error) {
	td, err := uc.jwt.Use(t)
	if err != nil {
		return td, fmt.Errorf("AuthUseCase - ValidateToken - uc.jwt.Use: %v", err)
	}
	return td, nil
}

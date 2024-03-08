package services

import (
	"context"
	"fmt"
	"github.com/qPyth/mobydev-internship-auth/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userStorage UserStorage
}

type UserStorage interface {
	CreateUser(ctx context.Context, email string, hashPass []byte) error
	GetUser(ctx context.Context, userID uint) (domain.User, error)
}

func NewAuthService(userStorage UserStorage) *AuthService {
	return &AuthService{userStorage: userStorage}
}

func (a *AuthService) SignUp(ctx context.Context, email, password string) error {
	//TODO:add validation
	op := "AuthService.SignUp"
	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s: bcrypt.GenerateFromPassword: %w", op, err)
	}
	return a.userStorage.CreateUser(ctx, email, hashPass)
}

func (a *AuthService) SignIn(ctx context.Context, email, password string) (token string, err error) {

}

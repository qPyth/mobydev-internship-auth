package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/qPyth/mobydev-internship-auth/internal/domain"
	"github.com/qPyth/mobydev-internship-auth/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userStorage  UserStorage
	TokenManager auth.TokenManager
}

type UserStorage interface {
	CreateUser(ctx context.Context, email string, hashPass []byte) error
	GetUser(ctx context.Context, email string) (domain.User, error)
	UpdateUser(ctx context.Context, req *domain.UserProfileUpdateReq) error
}

// NewUserService creates a new user service
func NewUserService(userStorage UserStorage, tokenManager auth.TokenManager) *UserService {
	return &UserService{userStorage: userStorage, TokenManager: tokenManager}
}

// SignUp creates a new user, returns domain.ErrEmailExists if user with such email already exists
func (u *UserService) SignUp(ctx context.Context, email, password string) error {
	op := "AuthService.SignUp"
	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s: bcrypt.GenerateFromPassword: %w", op, err)
	}
	return u.userStorage.CreateUser(ctx, email, hashPass)
}

// SignIn GetUser returns token for user by credentials. Returns domain.ErrInvalidCredentials if email or password is incorrect
func (u *UserService) SignIn(ctx context.Context, email, password string) (token string, err error) {
	op := "AuthService.SignIn"
	user, err := u.userStorage.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", domain.ErrInvalidCredentials
		}
		return "", fmt.Errorf("%s: userStorage.GetUser: %w", op, err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashPass), []byte(password)); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	token, err = u.TokenManager.NewJWT(user.ID)
	if err != nil {
		return "", fmt.Errorf("%s: tokenManager.NewJWT: %w", op, err)
	}

	return token, nil
}

// UpdateUserProfile updates user profile. Returns domain.ErrUserNotFound if user with such id not found
func (u *UserService) UpdateUserProfile(ctx context.Context, req domain.UserProfileUpdateReq) error {
	userID, ok := ctx.Value("userID").(float64)
	if !ok {
		return fmt.Errorf("userID not found in context")
	}
	req.ID = uint(userID)
	return u.userStorage.UpdateUser(ctx, &req)
}

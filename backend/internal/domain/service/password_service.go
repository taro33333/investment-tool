package service

import (
	"moneyget/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type PasswordService interface {
	HashPassword(password string) (string, error)
	ComparePasswords(hashedPassword, plainPassword string) error
	ValidatePassword(password string) error
}

type passwordService struct {
	minLength int
}

func NewPasswordService() PasswordService {
	return &passwordService{
		minLength: 8,
	}
}

func (s *passwordService) HashPassword(password string) (string, error) {
	if err := s.ValidatePassword(password); err != nil {
		return "", err
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (s *passwordService) ComparePasswords(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

func (s *passwordService) ValidatePassword(password string) error {
	if len(password) < s.minLength {
		return domain.ErrInvalidPassword
	}
	return nil
}

package usecase

import (
	"moneyget/internal/domain"
	"moneyget/internal/domain/service"
)

type UserUsecase interface {
	Register(name, email, password string) error
	Login(email, password string) (*domain.User, error)
	GetUserByID(id uint) (*domain.User, error)
}

type userUsecase struct {
	userRepo        domain.UserRepository
	passwordService service.PasswordService
}

func NewUserUsecase(repo domain.UserRepository, passwordService service.PasswordService) UserUsecase {
	return &userUsecase{
		userRepo:        repo,
		passwordService: passwordService,
	}
}

func (u *userUsecase) Register(name, email, password string) error {
	hashedPassword, err := u.passwordService.HashPassword(password)
	if err != nil {
		return err
	}

	user := &domain.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}
	return u.userRepo.Create(user)
}

func (u *userUsecase) Login(email, password string) (*domain.User, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if err := u.passwordService.ComparePasswords(user.Password, password); err != nil {
		return nil, domain.ErrInvalidCredentials
	}
	return user, nil
}

func (u *userUsecase) GetUserByID(id uint) (*domain.User, error) {
	return u.userRepo.FindByID(id)
}

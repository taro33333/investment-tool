package service

import (
	"time"

	"moneyget/internal/domain/constants"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateToken(userID uint) (string, error)
	ValidateToken(tokenString string) (uint, error)
}

type jwtService struct {
	secretKey []byte
}

func NewJWTService(secretKey string) JWTService {
	return &jwtService{
		secretKey: []byte(secretKey),
	}
}

func (s *jwtService) GenerateToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(constants.JWTTokenExpiration).Unix(),
	})
	return token.SignedString(s.secretKey)
}

func (s *jwtService) ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["user_id"].(float64))
		return userID, nil
	}
	return 0, jwt.ErrSignatureInvalid
}

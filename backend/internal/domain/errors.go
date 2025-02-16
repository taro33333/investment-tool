package domain

import (
	"errors"
	"fmt"
)

type DomainError struct {
	Code    string
	Message string
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

var (
	ErrInvalidEmail       = &DomainError{Code: "INVALID_EMAIL", Message: "Invalid email format"}
	ErrInvalidPassword    = &DomainError{Code: "INVALID_PASSWORD", Message: "Password must be at least 8 characters"}
	ErrUserNotFound       = &DomainError{Code: "USER_NOT_FOUND", Message: "User not found"}
	ErrInvalidCredentials = &DomainError{Code: "INVALID_CREDENTIALS", Message: "Invalid credentials"}
)

// 一般的なエラー
var (
	ErrNotFound     = errors.New("resource not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrUnauthorized = errors.New("unauthorized")
	ErrInternal     = errors.New("internal error")
)

// 投資関連のエラー
var (
	ErrInvalidInvestmentAmount = &DomainError{
		Code:    "INVALID_INVESTMENT_AMOUNT",
		Message: "investment amount is invalid",
	}

	ErrInvalidInvestmentType = &DomainError{
		Code:    "INVALID_INVESTMENT_TYPE",
		Message: "investment type is invalid",
	}

	ErrInvalidInvestmentStrategy = &DomainError{
		Code:    "INVALID_INVESTMENT_STRATEGY",
		Message: "investment strategy is invalid",
	}

	ErrPortfolioLimitExceeded = &DomainError{
		Code:    "PORTFOLIO_LIMIT_EXCEEDED",
		Message: "portfolio total amount would exceed maximum limit",
	}

	ErrInvestmentNotFound = &DomainError{
		Code:    "INVESTMENT_NOT_FOUND",
		Message: "investment not found in portfolio",
	}

	ErrAggressiveInvestmentLimitExceeded = &DomainError{
		Code:    "AGGRESSIVE_INVESTMENT_LIMIT_EXCEEDED",
		Message: "aggressive investments exceed maximum allowed ratio",
	}

	ErrDuplicateInvestment = &DomainError{
		Code:    "DUPLICATE_INVESTMENT",
		Message: "investment already exists in portfolio",
	}
)

// ポートフォリオ関連のエラー
var (
	ErrPortfolioNotFound    = errors.New("portfolio not found")
	ErrInvalidPortfolioData = errors.New("invalid portfolio data")
)

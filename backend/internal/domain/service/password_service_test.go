package service

import (
	"moneyget/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordService(t *testing.T) {
	passwordService := NewPasswordService()

	t.Run("ValidatePassword success", func(t *testing.T) {
		err := passwordService.ValidatePassword("validpass123")
		assert.NoError(t, err)
	})

	t.Run("ValidatePassword failure - too short", func(t *testing.T) {
		err := passwordService.ValidatePassword("short")
		assert.ErrorIs(t, err, domain.ErrInvalidPassword)
	})

	t.Run("HashPassword and ComparePasswords success", func(t *testing.T) {
		password := "validpass123"

		// ハッシュ化のテスト
		hashedPassword, err := passwordService.HashPassword(password)
		assert.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
		assert.NotEqual(t, password, hashedPassword)

		// パスワード比較のテスト
		err = passwordService.ComparePasswords(hashedPassword, password)
		assert.NoError(t, err)
	})

	t.Run("ComparePasswords failure - wrong password", func(t *testing.T) {
		password := "validpass123"
		wrongPassword := "wrongpass123"

		hashedPassword, _ := passwordService.HashPassword(password)
		err := passwordService.ComparePasswords(hashedPassword, wrongPassword)
		assert.Error(t, err)
	})

	t.Run("HashPassword failure - invalid password", func(t *testing.T) {
		_, err := passwordService.HashPassword("short")
		assert.ErrorIs(t, err, domain.ErrInvalidPassword)
	})
}

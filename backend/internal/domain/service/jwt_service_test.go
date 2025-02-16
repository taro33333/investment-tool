package service

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWTService(t *testing.T) {
	secretKey := "test-secret-key"
	jwtService := NewJWTService(secretKey)

	t.Run("GenerateToken and ValidateToken success", func(t *testing.T) {
		userID := uint(123)

		// トークン生成のテスト
		token, err := jwtService.GenerateToken(userID)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// トークン検証のテスト
		validatedUserID, err := jwtService.ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, userID, validatedUserID)
	})

	t.Run("ValidateToken with invalid token", func(t *testing.T) {
		// 不正なトークンのテスト
		_, err := jwtService.ValidateToken("invalid-token")
		assert.Error(t, err)
	})

	t.Run("ValidateToken with expired token", func(t *testing.T) {
		// 期限切れトークンの作成
		userID := uint(123)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": userID,
			"exp":     time.Now().Add(-time.Hour).Unix(), // 1時間前の期限切れトークン
		})
		expiredToken, _ := token.SignedString([]byte(secretKey))

		// 期限切れトークンの検証
		_, err := jwtService.ValidateToken(expiredToken)
		assert.Error(t, err)
	})
}

package sqlite

import (
	"moneyget/internal/domain"
	"testing"
	"time"
)

func TestUserRepository(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)

	// テストユーザーの作成
	user := &domain.User{
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
	}

	// Create のテスト
	t.Run("Create", func(t *testing.T) {
		err := repo.Create(user)
		if err != nil {
			t.Errorf("Failed to create user: %v", err)
		}
		if user.ID == "" {
			t.Error("Expected user ID to be set after creation")
		}
	})

	// FindByID のテスト
	t.Run("FindByID", func(t *testing.T) {
		found, err := repo.FindByID(user.ID)
		if err != nil {
			t.Errorf("Failed to find user by ID: %v", err)
		}
		if found.ID != user.ID {
			t.Errorf("Expected ID %s, got %s", user.ID, found.ID)
		}
		if found.Name != user.Name {
			t.Errorf("Expected Name %s, got %s", user.Name, found.Name)
		}
		if found.Email != user.Email {
			t.Errorf("Expected Email %s, got %s", user.Email, found.Email)
		}
	})

	// FindByEmail のテスト
	t.Run("FindByEmail", func(t *testing.T) {
		found, err := repo.FindByEmail(user.Email)
		if err != nil {
			t.Errorf("Failed to find user by email: %v", err)
		}
		if found.ID != user.ID {
			t.Errorf("Expected ID %s, got %s", user.ID, found.ID)
		}
		if found.Email != user.Email {
			t.Errorf("Expected Email %s, got %s", user.Email, found.Email)
		}
	})

	// Update のテスト
	t.Run("Update", func(t *testing.T) {
		user.Name = "Updated Name"
		err := repo.Update(user)
		if err != nil {
			t.Errorf("Failed to update user: %v", err)
		}

		found, err := repo.FindByID(user.ID)
		if err != nil {
			t.Errorf("Failed to find updated user: %v", err)
		}
		if found.Name != "Updated Name" {
			t.Errorf("Expected updated name 'Updated Name', got %s", found.Name)
		}
	})

	// Delete のテスト
	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(user.ID)
		if err != nil {
			t.Errorf("Failed to delete user: %v", err)
		}

		// 削除されたことを確認
		_, err = repo.FindByID(user.ID)
		if err == nil {
			t.Error("Expected error when finding deleted user")
		}
	})

	// 重複メールアドレスのテスト
	t.Run("DuplicateEmail", func(t *testing.T) {
		user1 := &domain.User{
			Name:      "User 1",
			Email:     "duplicate@example.com",
			Password:  "password1",
			CreatedAt: time.Now(),
		}
		user2 := &domain.User{
			Name:      "User 2",
			Email:     "duplicate@example.com",
			Password:  "password2",
			CreatedAt: time.Now(),
		}

		err := repo.Create(user1)
		if err != nil {
			t.Errorf("Failed to create first user: %v", err)
		}

		// 同じメールアドレスで2人目のユーザーを作成
		err = repo.Create(user2)
		if err == nil {
			t.Error("Expected error when creating user with duplicate email")
		}
	})
}

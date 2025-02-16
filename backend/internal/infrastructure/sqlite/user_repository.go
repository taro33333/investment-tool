package sqlite

import (
	"database/sql"
	"moneyget/internal/domain"

	_ "github.com/mattn/go-sqlite3"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	// テーブルの作成
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		)
	`)
	if err != nil {
		panic(err)
	}
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	query := `INSERT INTO users (name, email, password) VALUES (?, ?, ?)`
	result, err := r.db.Exec(query, user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = uint(id)
	return nil
}

func (r *userRepository) FindByID(id uint) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, name, email, password FROM users WHERE id = ?`
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, name, email, password FROM users WHERE email = ?`
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	query := `UPDATE users SET name = ?, email = ?, password = ? WHERE id = ?`
	_, err := r.db.Exec(query, user.Name, user.Email, user.Password, user.ID)
	return err
}

func (r *userRepository) Delete(id uint) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

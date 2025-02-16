package sqlite

import (
	"database/sql"
	"moneyget/internal/domain"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	user.ID = uuid.New().String()
	query := `INSERT INTO users (id, name, email, password, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, user.ID, user.Name, user.Email, user.Password, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) FindByID(id string) (*domain.User, error) {
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
	result, err := r.db.Exec(query, user.Name, user.Email, user.Password, user.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *userRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

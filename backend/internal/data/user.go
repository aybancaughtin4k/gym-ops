package data

import (
	"context"
	"errors"
	"time"

	"github.com/aybancaughtin4k/gymops/backend/internal/validator"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("email already exist")
)

type UserModel struct {
	dbpool *pgxpool.Pool
}

type User struct {
	ID        int64     `json:"id"`
	Fullname  string    `json:"fullname"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type password struct {
	plaintext *string
	hash      []byte
}

func (m UserModel) GetUserByEmail(email string) (*User, error) {
	var user User

	query := `SELECT * FROM users WHERE email = @email`
	args := pgx.NamedArgs{
		"email": email,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.dbpool.QueryRow(ctx, query, args).Scan(&user.ID, &user.Fullname, &user.Email, &user.Password.hash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) GetUserByUsername(username string) (*User, error) {
	var user User

	query := `SELECT * FROM users WHERE username = @username`
	args := pgx.NamedArgs{
		"username": username,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.dbpool.QueryRow(ctx, query, args).Scan(&user.ID, &user.Fullname, &user.Email, &user.Password.hash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) Insert(user *User) error {
	query := `INSERT INTO users (fullname, email, password)
		VALUES (@fullname, @email, @password)
		RETURNING id, created_at, updated_at`
	args := pgx.NamedArgs{
		"fullname": user.Fullname,
		"email":    user.Email,
		"password": user.Password.hash,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.dbpool.QueryRow(ctx, query, args).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		switch {
		case err.Error() == `ERROR: duplicate key value violates unique constraint "users_email_key" (SQLSTATE 23505)`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) Update(userId int64, user *User) error {
	query := `UPDATE users
	SET fullname = @fullname, email = @email, password = @password, updated_at = @updated_at
	WHERE id = @userId
	RETURNING fullname, email, password, updated_at`
	args := pgx.NamedArgs{
		"fullname":   user.Fullname,
		"email":      user.Email,
		"password":   user.Password,
		"updated_at": time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return m.dbpool.QueryRow(ctx, query, args).Scan(&user.Fullname, &user.Email, &user.Password, &user.UpdatedAt)
}

func (m UserModel) Delete(userId int64) error {
	query := `DELETE FROM users WHERE id = @userId`
	args := pgx.NamedArgs{
		"userId": userId,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := m.dbpool.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	row := result.RowsAffected()
	if row == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// Hash user's password
func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

// Check if user matches the hashed password
func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "email is required")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "invalid email format")
}

func ValidatePasswordPlainText(v *validator.Validator, password string) {
	v.Check(password != "", "password", "password is required")
	v.Check(len(password) >= 8, "password", "password must be greater than 8 characters")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Fullname != "", "fullname", "fullname is required")
	v.Check(len(user.Fullname) >= 10, "fullname", "fullname must be greater than 10 characters")

	ValidateEmail(v, user.Email)

	if user.Password.plaintext != nil {
		ValidatePasswordPlainText(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"github.com/qPyth/mobydev-internship-auth/internal/domain"
	"strings"
	"time"
)

var (
	migrationsPath = "file://migrations"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) *Storage {
	const op = "sqlite.New"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		panic("failed to open db: " + err.Error())
	}
	err = migrateDB(db)
	return &Storage{db: db}
}

func migrateDB(db *sql.DB) error {
	op := "sqlite.migrate"
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("%s: sqlite.WithInstance: %w", op, err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"sqlite3", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

// CreateUser creates a new user, returns domain.ErrEmailExists if user with such email already exists
func (s *Storage) CreateUser(ctx context.Context, email string, hashPass []byte) error {
	op := "sqlite.CreateUser"
	now := time.Now()
	stmt, err := s.db.Prepare("INSERT INTO users(email, password, created_at, updated_at) VALUES(?, ?, ?,?)")
	if err != nil {
		return fmt.Errorf("%s: db.Prepare: %w", op, err)
	}
	_, err = stmt.ExecContext(ctx, email, string(hashPass), now, now)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return domain.ErrEmailExists
		}
		return fmt.Errorf("%s: stmt.Exec: %w", op, err)
	}
	return nil
}

// GetUser returns user by id. Returns domain.ErrUserNotFound if user not found
func (s *Storage) GetUser(ctx context.Context, email string) (domain.User, error) {
	op := "sqlite.GetUser"

	var user domain.User

	row := s.db.QueryRowContext(ctx, "SELECT id, email, password FROM users WHERE email = ?", email)
	err := row.Scan(&user.ID, &user.Email, &user.HashPass)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, domain.ErrUserNotFound
		}
		return user, fmt.Errorf("%s: row.Scan: %w", op, err)
	}
	return user, nil
}

// UpdateUser updates user profile. Returns domain.ErrUserNotFound if user not found
func (s *Storage) UpdateUser(ctx context.Context, update *domain.UserProfileUpdateReq) error {
	op := "sqlite.UpdateUser"
	row := s.db.QueryRowContext(ctx, "SELECT id FROM users WHERE id = ?", update.ID)
	var id uint
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, domain.ErrUserNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	var queryBuilder strings.Builder
	var args []interface{}

	queryBuilder.WriteString("UPDATE users SET ")

	if update.Name != nil {
		queryBuilder.WriteString("name = ?, ")
		args = append(args, *update.Name)
	}
	if update.Email != nil {
		queryBuilder.WriteString("email = ?, ")
		args = append(args, *update.Email)
	}

	if update.BDay != nil {
		queryBuilder.WriteString("b_day = ?, ")
		args = append(args, *update.BDay)
	}
	if update.PhoneNumber != nil {
		queryBuilder.WriteString("phone_number = ?, ")
		args = append(args, *update.PhoneNumber)
	}

	queryBuilder.WriteString("updated_at = ?, ")
	args = append(args, time.Now())

	query := strings.TrimSuffix(queryBuilder.String(), ", ")

	query += "WHERE id = ?"
	args = append(args, update.ID)

	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

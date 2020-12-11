// Package interestvalue provides struct values and associated operations for user handling including persistence in postgres
package uservalue

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// User represents someone with access to our system.
type User struct {
	ID           string         `Db:"user_id" json:"id"`
	Name         string         `Db:"name" json:"name"`
	Email        string         `Db:"email" json:"email"`
	Roles        pq.StringArray `Db:"roles" json:"roles"`
	PasswordHash []byte         `Db:"password_hash" json:"-"`
	DateCreated  time.Time      `Db:"date_created" json:"date_created"`
	DateUpdated  time.Time      `Db:"date_updated" json:"date_updated"`
}

// NewUser contains information needed to create a new User.
type NewUser struct {
	Name            string   `json:"name" validate:"required"`
	Email           string   `json:"email" validate:"required"`
	Roles           []string `json:"roles" validate:"required"`
	Password        string   `json:"password" validate:"required"`
	PasswordConfirm string   `json:"password_confirm" validate:"eqfield=Password"`
}

func AddUser(ctx context.Context, db *sqlx.DB, n NewUser, now time.Time) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(n.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "generating password hash")
	}

	u := User{
		ID:           uuid.New().String(),
		Name:         n.Name,
		Email:        n.Email,
		PasswordHash: hash,
		Roles:        n.Roles,
		DateCreated:  now.UTC(),
		DateUpdated:  now.UTC(),
	}

	const q = `INSERT INTO users
		(user_id, name, email, password_hash, roles, date_created, date_updated)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = db.ExecContext(
		ctx, q,
		u.ID, u.Name, u.Email,
		u.PasswordHash, u.Roles,
		u.DateCreated, u.DateUpdated,
	)
	if err != nil {
		fmt.Printf("\ndb err is %T %v", err, err)
		return nil, errors.Wrap(err, "inserting usermgmt")
	}

	return &u, nil
}

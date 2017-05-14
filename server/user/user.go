package user

import (
    "github.com/pkg/errors"
    "database/sql"
    "github.com/jmoiron/sqlx"
)

var (
    ErrPasswordConfirmationDoesNotMatch = errors.New("Password confirmation doesn't match.")
    ErrUsernameAlreadyTaken = errors.New("Username already taken")
)
type User interface {
    Create(username string, password string, passwordConfirmation string) error
}



type Builder struct {
    DB *sql.DB
}

func (builder Builder) Build() User {
    return &user{
        db: sqlx.NewDb(builder.DB, ""),
    }
}

type user struct {
    db *sqlx.DB
}

func (u *user) Create(username string, password string, passwordConfirmation string) error {
    return nil
}
package session

import (
    "github.com/pkg/errors"
    "database/sql"
    "github.com/jmoiron/sqlx"
)

var (
    ErrBadCredentials = errors.New("Bad credentials")
    ErrTokenNotFound = errors.New("Token not found")
)

type Data struct {
    Username string
}

type Session interface {
    Create(username string, password string) (token string, err error)
    Delete(token string) error
    Valid(token string) error
    Retrieve(token string) (*Data, error)
}

type Builder struct {
    DB *sql.DB
}

func (builder Builder) Build() Session {
    return &session{
        db: sqlx.NewDb(builder.DB, ""),
    }
}

type session struct {
    db *sqlx.DB
}

func (s *session) Create(username string, password string) (token string, err error) {
    return "", nil
}

func (s *session) Delete(token string) error {
    return nil
}

func (s *session) Valid(token string) error {
    return nil
}

func (s *session) Retrieve(token string) (*Data, error) {
    return nil, nil
}


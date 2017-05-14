package session

import (
    "github.com/pkg/errors"
    "database/sql"
    "github.com/jmoiron/sqlx"
    "golang.org/x/crypto/bcrypt"
    "github.com/satori/go.uuid"
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

func (s *session) Create(username string, password string) (string, error) {
    var (
        hashedPassword string
        userID int64
    )
    const selQuery = `SELECT password, user_id FROM "user" WHERE username = $1;`
    if err := s.db.QueryRow(selQuery, username).Scan(&hashedPassword, &userID); err == sql.ErrNoRows {
        return "", ErrBadCredentials
    } else if err != nil {
        return "", err
    } else  if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
        return "", ErrBadCredentials
    }

    key := uuid.NewV4().String()
    const insQuery = `INSERT INTO "user_session"("session_key", "user_id") VALUES($1, $2);`
    if _, err := s.db.Exec(insQuery, key, userID); err != nil {
        return "", err
    }

    return key, nil
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


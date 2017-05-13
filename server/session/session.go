package session

import "github.com/pkg/errors"

var (
    ErrBadCredentials = errors.New("Bad credentials")
    ErrTokenNotFound = errors.New("Token not found")
    ErrInvalidToken = errors.New("Invalid Token.")
)
type Session interface {
    Create(username string, password string) (token string,  err error)
    Delete(token string) error
    Valid(token string) error
}


func New() Session {
    return Session(nil)
}
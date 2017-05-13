package session

import "github.com/pkg/errors"

var (
    ErrBadCredentials = errors.New("Bad credentials")
    ErrTokenNotFound = errors.New("Token not found")
)
type Session interface {
    Create(username string, password string) (key string,  err error)
    Delete(key string) error
}


func New() Session {
    return Session(nil)
}
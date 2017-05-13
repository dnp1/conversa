package session

import "github.com/pkg/errors"

var (
    ErrBadCredentials = errors.New("Bad credentials")
)
type Session interface {
    Create(username string, password string) (key string,  err error)
}


func New() Session {
    return Session(nil)
}
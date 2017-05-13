package session

import "github.com/pkg/errors"

var (
    ErrBadCredentials = errors.New("Bad credentials")
    ErrTokenNotFound = errors.New("Token not found")
)

type Data struct {
    Username string
}

type Session interface {
    Create(username string, password string) (token string,  err error)
    Delete(token string) error
    Valid(token string) error
    Retrieve(token string) (*Data, error)
}


func New() Session {
    return Session(nil)
}
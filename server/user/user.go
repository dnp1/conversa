package user

import "github.com/pkg/errors"

var (
    ErrPasswordConfirmationDoesNotMatch = errors.New("Password confirmation doesn't match.")
    ErrUsernameAlreadyTaken = errors.New("Username already taken")
)
type User interface {
    Create(username string, password string, passwordConfirmation string) error
}


func New() User {
    return User(nil)
}
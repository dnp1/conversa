package errors

import "github.com/pkg/errors"

type Error interface {
    error
    Validation() bool
    Conflict() bool
    Internal() bool
    Empty() bool
    Authorization() bool
    Authentication() bool
}

type myError struct {
    error
    validation     bool
    conflict       bool
    internal       bool
    empty          bool
    authorization  bool
    authentication bool
}

func (err *myError) Validation() bool {
    return err.validation
}
func (err *myError) Conflict() bool {
    return err.conflict
}
func (err *myError) Internal() bool {
    return err.internal
}
func (err *myError) Empty() bool {
    return err.empty
}
func (err *myError) Authorization() bool {
    return err.authorization
}
func (err *myError) Authentication() bool {
    return err.authentication
}

func FromString(msg string) error {
    return errors.New(msg)
}

func Validation(err error) *myError {
    return &myError{error: err, validation: true}
}
func Conflict(err error) *myError {
    return &myError{error: err, conflict: true}
}
func Internal(err error) *myError {
    return &myError{error: err, internal: true}
}
func Empty(err error) *myError {
    return &myError{error: err, empty: true}
}
func Authorization(err error) *myError {
    return &myError{error: err, authorization: true}
}
func Authentication(err error) *myError {
    return &myError{error: err, authentication: true}
}

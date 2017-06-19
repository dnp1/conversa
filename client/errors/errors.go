package errors

import (
    "errors"
    "net/http"
    "fmt"
)

type myError struct {
    error
    badRequest     bool
    conflict       bool
    server         bool
    notFound       bool
    authorization  bool
    authentication bool
    unexpected     bool
}

func (err *myError) BadRequest() bool {
    return err.badRequest
}

func (err *myError) Conflict() bool {
    return err.conflict
}

func (err *myError) Server() bool {
    return err.server
}

func (err *myError) NotFound() bool {
    return err.notFound
}

func (err *myError) Authorization() bool {
    return err.authorization
}
func (err *myError) Authentication() bool {
    return err.authentication
}

func (err *myError) Unexpected() bool {
    return err.unexpected
}

func FromString(msg string) error {
    return errors.New(msg)
}

func FromHttpStatus(code int, msg string) *myError {
    err := fmt.Errorf("Message:%q\nCode:%d\n", msg, code)
    switch {
    case code >=200 && code <300:
        return nil
    case code >= 500 && code < 600:
        return Server(err)
    case code == http.StatusConflict:
        return Conflict(err)
    case code == http.StatusNotFound:
        return NotFound(err)
    case code == http.StatusUnauthorized:
        return Authentication(err)
    case code == http.StatusForbidden:
        return Authorization(err)
    case code == http.StatusBadRequest:
        return BadRequest(err)
    default:
        return Unexpected(err)
    }
}

func BadRequest(err error) *myError {
    return &myError{error: err, badRequest: true}
}
func Conflict(err error) *myError {
    return &myError{error: err, conflict: true}
}
func Server(err error) *myError {
    return &myError{error: err, server: true}
}
func NotFound(err error) *myError {
    return &myError{error:err, notFound: true}
}
func Authorization(err error) *myError {
    return &myError{error: err, authorization: true}
}
func Authentication(err error) *myError {
    return &myError{error: err, authentication: true}
}
func Unexpected(err error) *myError {
    return &myError{error: err, unexpected: true}
}

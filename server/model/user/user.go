package user

import (
    "database/sql"
    "golang.org/x/crypto/bcrypt"
    "fmt"
    "unicode/utf8"
    "regexp"
    "github.com/dnp1/conversa/server/errors"
)

const (
    UsernameMinLength = 3
    UsernameMaxLength = 255
    PasswordMinLength = 8
)

var (
    //errors
    ErrPasswordConfirmationDoesNotMatch = errors.Validation(
        errors.FromString("Password confirmation doesn't match."))

    ErrUsernameAlreadyTaken = errors.Conflict(errors.FromString("Username already taken"))
    ErrUsernameWrongLength  = errors.Validation(fmt.Errorf(
        "Username length is valid, mininum is %d and maximum is %d",
        UsernameMinLength,
        UsernameMaxLength,
    ))
    ErrPasswordTooShort = errors.Validation(
        fmt.Errorf("Passord must be longer than %d characters", PasswordMinLength))
    ErrUsernameHasInvalidCharacters = errors.Validation(
        errors.FromString("Username can only contain alphanumeric characters and underscores."))
)

func New(db *sql.DB, bCryptCost int) *model {
    if bCryptCost == 0 {
        bCryptCost = bcrypt.DefaultCost
    }
    return &model{
        db:         db,
        bCryptCost: bCryptCost,
    }
}

type model struct {
    db         *sql.DB
    bCryptCost int
}

var regexpUsername = regexp.MustCompile("^[a-zA-Z0-9_]+$")

func Validate(username, password, passwordConfirmation string) errors.Error {
    if password != passwordConfirmation {
        return ErrPasswordConfirmationDoesNotMatch
    } else
    if length := utf8.RuneCountInString(username); length < UsernameMinLength || length > UsernameMaxLength {
        return ErrUsernameWrongLength
    }
    if length := utf8.RuneCountInString(password); length < PasswordMinLength {
        return ErrPasswordTooShort
    }
    if hasInvalidChars := !regexpUsername.MatchString(username); hasInvalidChars {
        return ErrUsernameHasInvalidCharacters
    }

    return nil
}

func (u *model) Create(username string, password string, passwordConfirmation string) errors.Error {
    var id int64;
    const query = `INSERT INTO "model"("username", "password") VALUES($1, $2)
        ON CONFLICT ON CONSTRAINT "uq_username" DO NOTHING RETURNING id;`
    if err := Validate(username, password, passwordConfirmation); err != nil {
        return err
    } else if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), u.bCryptCost); err != nil {
        return errors.Internal(err)
    } else if err := u.db.QueryRow(query, username, string(hashedPassword)).Scan(&id); err == sql.ErrNoRows {
        return ErrUsernameAlreadyTaken
    } else if err != nil {
        return errors.Internal(err)
    }
    return nil
}

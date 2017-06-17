package user

import (
    "github.com/pkg/errors"
    "database/sql"
    "golang.org/x/crypto/bcrypt"
    "fmt"
    "unicode/utf8"
    "regexp"
)

const (
    UsernameMinLength = 3
    UsernameMaxLength = 255
    PasswordMinLength = 8
)

var (
    //errors
    ErrPasswordConfirmationDoesNotMatch = errors.New("Password confirmation doesn't match.")
    ErrUsernameAlreadyTaken = errors.New("Username already taken")
    ErrUsernameWrongLength = fmt.Errorf(
        "Username length is valid, mininum is %d and maximum is %d",
        UsernameMinLength,
        UsernameMaxLength,
    )
    ErrPasswordTooShort = fmt.Errorf("Passord must be longer than %d characters", PasswordMinLength)
    ErrUsernameHasInvalidCharacters = errors.New("Username can only contain alphanumeric characters and underscores.")
)

type User interface {
    Create(username string, password string, passwordConfirmation string) error
}

type Builder struct {
    DB         *sql.DB
    BCryptCost int
}

func (builder Builder) Build() User {
    if builder.BCryptCost == 0 {
        builder.BCryptCost = bcrypt.DefaultCost
    }
    return &user{
        db: builder.DB,
        bCryptCost: builder.BCryptCost,
    }
}

type user struct {
    db         *sql.DB
    bCryptCost int
}

var regexpUsername = regexp.MustCompile("^[a-zA-Z0-9_]+$")

func Validate(username, password, passwordConfirmation string) error {
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

func (u *user) Create(username string, password string, passwordConfirmation string) error {
    var id int64;
    const query = `INSERT INTO "user"("username", "password") VALUES($1, $2)
        ON CONFLICT ON CONSTRAINT "uq_username" DO NOTHING RETURNING id;`
    if err := Validate(username, password, passwordConfirmation); err != nil {
        return err
    } else if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), u.bCryptCost); err != nil {
        return err
    } else if err := u.db.QueryRow(query, username, string(hashedPassword)).Scan(&id); err == sql.ErrNoRows {
        return ErrUsernameAlreadyTaken
    } else if err != nil {
        return err
    }

    return nil
}
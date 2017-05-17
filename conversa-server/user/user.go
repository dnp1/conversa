package user

import (
    "github.com/pkg/errors"
    "database/sql"
    "github.com/jmoiron/sqlx"
    "golang.org/x/crypto/bcrypt"
)

var (
    ErrPasswordConfirmationDoesNotMatch = errors.New("Password confirmation doesn't match.")
    ErrUsernameAlreadyTaken = errors.New("Username already taken")
)
type User interface {
    Create(username string, password string, passwordConfirmation string) error
}

type Builder struct {
    DB *sql.DB
    BCryptCost int
}

func (builder Builder) Build() User {
    if builder.BCryptCost == 0 {
        builder.BCryptCost = bcrypt.DefaultCost
    }
    return &user{
        db: sqlx.NewDb(builder.DB, ""),
        bCryptCost: builder.BCryptCost ,
    }
}

type user struct {
    db *sqlx.DB
    bCryptCost int
}

func (u *user) Create(username string, password string, passwordConfirmation string) error {
    //TODO:validate username with a regexp
    if password != passwordConfirmation {
        return ErrPasswordConfirmationDoesNotMatch
    }
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), u.bCryptCost)
    if err != nil {
        return err
    }
    const query = `INSERT INTO "user"("username", "password") VALUES($1, $2)
        ON CONFLICT ON CONSTRAINT "uq_username" DO NOTHING RETURNING id;`
    var id int64;
    if err := u.db.QueryRow(query, username, string(hashedPassword)).Scan(&id); err == sql.ErrNoRows {
        return ErrUsernameAlreadyTaken
    } else if err != nil {
        return err
    }
    return nil
}
package room

import (
    "github.com/pkg/errors"
    "database/sql"
    "github.com/jmoiron/sqlx"
)

var (
    ErrCouldNotInsert = errors.New("Couldn't insert!")
    ErrCouldNotDelete = errors.New("Couldn't delete!")
    ErrCouldNotRename = errors.New("Couldn't rename!")
    ErrCouldNotRetrieveRooms = errors.New("Select has failed!")
)

type Data struct {
    Username string `json:"username"`
    Name     string `json:"name"`
}

type Room interface {
    Create(username string, name string) error
    Delete(username string, name string) error
    Rename(username, oldName, newName string) error
    All() ([]Data, error)
    AllByUser(username string) ([]Data, error)
}

type Builder struct {
    DB *sql.DB
}

func (builder Builder) Build() Room {
    return &room{
        db: sqlx.NewDb(builder.DB, ""),
    }
}

type room struct {
    db *sqlx.DB
}

func (r *room) Create(username string, name string) error {
    return nil
}

func (r *room) Delete(username string, name string) error {
    return nil
}

func (r *room) Rename(username, oldName, newName string) error {
    return nil
}

func (r *room) All() ([]Data, error) {
    return nil, nil
}

func (r *room) AllByUser(username string) ([]Data, error) {
    return nil, nil
}
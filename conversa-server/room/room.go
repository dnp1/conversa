package room

import (
    "github.com/pkg/errors"
    "database/sql"
    "log"
    "regexp"
    "fmt"
    "unicode/utf8"
)

const (
    RoomNameMinLength = 3
    RoomNameMaxLength = 255
)
var (
    ErrRoomNameAlreadyExists = errors.New("Couldn't insert!")
    ErrCouldNotDelete = errors.New("Couldn't delete!")
    ErrCouldNotRename = errors.New("Couldn't rename!")
    ErrCouldNotRetrieveRooms = errors.New("Select has failed!")
    ErrRoomNameWrongLength = fmt.Errorf(
        "Room's name length is valid, mininum is %d and maximum is %d",
        RoomNameMinLength,
        RoomNameMaxLength,
    )
    ErrRoomNameHasInvalidCharacters = errors.New("Room name can only contain alphanumeric characters and underscores.")
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
        db: builder.DB,
    }
}

type room struct {
    db *sql.DB
}

var regexpRoom = regexp.MustCompile("[a-zA-Z0-9_]")

func (r *room) Create(username string, name string) error {
    const query = `INSERT INTO "room"("name", "username", "user_id")
        SELECT $1, $2::TEXT, u.id FROM "user" u WHERE u."username" = $2
        ON CONFLICT ON CONSTRAINT "uq_name" DO NOTHING RETURNING id;
    `
    if length:=utf8.RuneCountInString(name); length < RoomNameMinLength || length > RoomNameMaxLength {
        return ErrRoomNameWrongLength
    }
    if hasInvalidChars := !regexpRoom.MatchString(name); hasInvalidChars {
        return ErrRoomNameHasInvalidCharacters
    }

    var id int64
    if err := r.db.QueryRow(query, name, username).Scan(&id); err == sql.ErrNoRows {
        return ErrRoomNameAlreadyExists
    } else if err != nil{
        return err
    }
    return nil
}

func (r *room) Delete(username string, name string) error {
    const query = `DELETE FROM "room" WHERE
        "name"=$1 AND "username" = $2
    `
    if _, err := r.db.Exec(query, name, username); err != nil {
        log.Println(err)
        return ErrCouldNotDelete
    }
    return nil
}

func (r *room) Rename(username, oldName, newName string) error {
    const query = `UPDATE "room" SET "name" = $3 WHERE "username"=$1 AND "name"=$2;`
    if hasInvalidChars := !regexpRoom.MatchString(newName); hasInvalidChars {
        return ErrRoomNameHasInvalidCharacters
    }
    if _, err := r.db.Exec(query, username, oldName, newName); err != nil {
        return err
    }
    return nil
}

func (r *room) All() ([]Data, error) {
    const query = `SELECT username, name FROM "room" ORDER BY "username", "name";`
    if rows, err := r.db.Query(query); err != nil {
        return nil, err
    } else {
        var set = make([]Data, 0)
        for rows.Next() {
            var row Data
            if err := rows.Scan(&row.Username, &row.Name); err != nil {
                return nil, err
            }
            set = append(set, row)
        }
        return set, nil
    }
}

func (r *room) AllByUser(username string) ([]Data, error) {
    const query = `SELECT name FROM "room" where "username"=$1 ORDER BY "name";`
    if rows, err := r.db.Query(query); err != nil {
        return nil, err
    } else {
        var set = make([]Data, 0)
        for rows.Next() {
            var row  = Data{Name:username}
            if err := rows.Scan(&row.Name); err != nil {
                return nil, err
            }
            set = append(set, row)
        }
        return set, nil
    }
}
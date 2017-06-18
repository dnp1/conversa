package room

import (
    "database/sql"
    "regexp"
    "fmt"
    "unicode/utf8"
    "github.com/dnp1/conversa/server/errors"
    "github.com/dnp1/conversa/server/data/room"
)

const (
    NameMinLength = 3
    NameMaxLength = 255
)

var (
    ErrRoomNameAlreadyExists = errors.Conflict(errors.FromString("Couldn't insert!"))
    ErrCouldNotDelete        = errors.FromString("Couldn't delete!")
    ErrCouldNotRename        = errors.FromString("Couldn't rename!")
    ErrCouldNotRetrieveRooms = errors.FromString("Select has failed!")
    ErrRoomNameWrongLength   = errors.Validation(fmt.Errorf(
        "Model's name length is invalid, mininum is %d and maximum is %d",
        NameMinLength,
        NameMaxLength,
    ))
    ErrRoomNameHasInvalidCharacters = errors.Validation(
        errors.FromString("Model name can only contain alphanumeric characters and underscores."))
)

func New(db *sql.DB) *model {
    return &model{
        db: db,
    }
}

type model struct {
    db *sql.DB
}

var regexpRoom = regexp.MustCompile("[a-zA-Z0-9_]")

func (r *model) Create(username string, name string) errors.Error {
    const query = `INSERT INTO "model"("name", "username", "user_id")
        SELECT $1, $2::TEXT, u.id FROM "user" u WHERE u."username" = $2
        ON CONFLICT ON CONSTRAINT "uq_name" DO NOTHING RETURNING id;
    `
    if length := utf8.RuneCountInString(name); length < NameMinLength || length > NameMaxLength {
        return ErrRoomNameWrongLength
    }
    if hasInvalidChars := !regexpRoom.MatchString(name); hasInvalidChars {
        return ErrRoomNameHasInvalidCharacters
    }

    var id int64
    if err := r.db.QueryRow(query, name, username).Scan(&id); err == sql.ErrNoRows {
        return ErrRoomNameAlreadyExists
    } else if err != nil {
        return errors.Internal(err)
    }
    return nil
}

func (r *model) Delete(username string, name string) errors.Error {
    const query = `DELETE FROM "model" WHERE
        "name"=$1 AND "username" = $2
    `
    if _, err := r.db.Exec(query, name, username); err != nil {
        return errors.Internal(err)
    }
    return nil
}

func (r *model) Rename(username, oldName, newName string) errors.Error {
    const query = `UPDATE "model" SET "name" = $3 WHERE "username"=$1 AND "name"=$2;`
    if hasInvalidChars := !regexpRoom.MatchString(newName); hasInvalidChars {
        return ErrRoomNameHasInvalidCharacters
    }
    if _, err := r.db.Exec(query, username, oldName, newName); err != nil {
        return errors.Internal(err)
    }
    return nil
}

func (r *model) All() ([]room.Data, errors.Error) {
    const query = `SELECT username, name FROM "model" ORDER BY "username", "name";`
    if rows, err := r.db.Query(query); err != nil {
        return nil, errors.Internal(err)
    } else {
        var set = make([]room.Data, 0)
        for rows.Next() {
            var row room.Data
            if err := rows.Scan(&row.Username, &row.Name); err != nil {
                return nil, errors.Internal(err)
            }
            set = append(set, row)
        }
        return set, nil
    }
}

package room

import "github.com/pkg/errors"

var ErrCouldNotInsert = errors.New("Couldn't insert!")

type Room interface{
    Create(username string, name string) error
}

func New() Room {
    return Room(nil)
}
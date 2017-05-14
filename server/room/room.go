package room

import "github.com/pkg/errors"

var (
    ErrCouldNotInsert = errors.New("Couldn't insert!")
    ErrCouldNotDelete = errors.New("Couldn't delete!")
    ErrCouldNotRename = errors.New("Couldn't rename!")
    ErrCouldNotRetrieveRooms = errors.New("Select has failed!")
)


type Data struct {
    Username string `json:"username"`
    Name string `json:"name"`
}

type Room interface{
    Create(username string, name string) error
    Delete(username string, name string) error
    Rename(username, oldName, newName string) error
    All() ([]Data, error)
    AllByUser(username string) ([]Data, error)
}

func New() Room {
    return Room(nil)
}
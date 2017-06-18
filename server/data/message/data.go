package message

import "time"

type Data struct {
    ID int64 `json:"id"`
    RoomID int64 `json:"roomId"`
    UserID int64 `json:"userId"`
    Username string `json:"username"`
    Content string `json:"content"`
    CreationDate time.Time
    EditionDate time.Time
}


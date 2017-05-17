package message

import (
    "time"
    "database/sql"
)

type Data struct {
    ID int64 `json:"id"`
    RoomID int64 `json:"roomId"`
    UserID int64 `json:"userId"`
    Username string `json:"username"`
    Content string `json:"content"`
    CreationDate time.Time
    EditionDate time.Time
}


type Message interface {
    Create(username , roomName, senderName, content string) error
    Edit(messageID int , content string) error
    Delete(messageID int) error
    All(username , roomName string, limit, offset int64) ([]Data, error)
}

type Builder struct {
    DB *sql.DB
}


func (builder Builder) Build() Message {
    return Message(nil)
}
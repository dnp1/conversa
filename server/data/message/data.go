package message

import "time"

type Data struct {
    ID string `json:"id"`
    RoomName string `json:"roomName"`
    RoomUsername string `json:"roomUsername"`
    OwnerUsername string `json:"ownerUsername"`
    Content string `json:"content"`
    CreationDate time.Time
    EditionDate time.Time
}

type EventData struct {
    Event string `json:"event"` //delete, edit,create
    Data
}


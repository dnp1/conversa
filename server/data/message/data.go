package message

import "time"

type Data struct {
    ID               string `json:"id"`
    RoomName         string `json:"roomName"`
    RoomUsername     string `json:"roomUsername"`
    OwnerUsername    string `json:"ownerUsername"`
    Content          string `json:"content"`
    CreationDatetime time.Time `json:"creationDatetime"`
    EditionDatetime  time.Time `json:"editionDatetime"`
}

type EventData struct {
    Event string `json:"event"` //delete, edit,create
    Data
}


package message

import (
    "database/sql"
    "strings"
    "github.com/dnp1/conversa/server/errors"
    "github.com/dnp1/conversa/server/data/message"
    "sync"
    "time"
)

var ( //errors
    ErrMessageIsEmpty = errors.Validation(errors.FromString("Message is empty, invalid!"))
)

func New(db *sql.DB) *model {
    return &model{
        db: db,
        listeners: make(map[mapKey][](chan *message.EventData), 1000),
    }
}

type mapKey struct {
    username string
    roomName string
}

type model struct {
    db        *sql.DB
    sync.RWMutex
    listeners map[mapKey][](chan *message.EventData)
}

func (m *model) Create(username, roomName, senderName, content string) errors.Error {
    const query = `INSERT INTO "message"("room_id", "user_id", "content")
    SELECT r.id, u.id, $4
        FROM room r
        LEFT JOIN "user" u ON u.username = $3
        WHERE r.username = $1 AND r.name = $2
    RETURNING id, creation_datetime, edition_datetime;`

    var id string
    var creationDatetime time.Time
    var editionDatetime time.Time
    content = strings.TrimSpace(content)
    if content == "" {
        return ErrMessageIsEmpty
    }

    if err := m.db.QueryRow(
        query,
        username,
        roomName,
        senderName,
        content,
    ).Scan(&id, &creationDatetime, &editionDatetime); err == sql.ErrNoRows {
        return errors.Validation(err)
    } else if err != nil {
        return errors.Internal(err)
    }

    go m.Broadcast(eventCreation, &message.Data{
        ID:               id,
        Content:          content,
        OwnerUsername:    senderName,
        RoomName:         roomName,
        RoomUsername:     username,
        CreationDatetime: creationDatetime,
        EditionDatetime:  editionDatetime,
    })

    return nil
}

func (m *model) Edit(username, roomName, messageOwner, messageID, content string) errors.Error {
    const query = `
    WITH (
        SELECT m.id as "msg_id", r.id as "room_id", u.id as "user_id"
                FROM model m
                INNER JOIN room r ON r.id = m.room_id
                INNER JOIN "user" u ON u.id = m.user_id
                WHERE
                    m.id = $4
                    u.username = $3
                    AND r.name = $1 AND r.username = $2
    ) AS "t"
    UPDATE "message" SET "content" = $5, "edition_datetime" = current_timestamp
        WHERE "message".id = t.msg_id
        AND message.room_id = t.room_id
        AND message.user_id = t.user_id
    RETURNING id, "creation_datetime", "edition_datetime"
;`

    var id string
    var creationDatetime time.Time
    var editionDatetime time.Time

    content = strings.TrimSpace(content)
    if content == "" {
        return ErrMessageIsEmpty
    }

    if err := m.db.QueryRow(
        query,
        username,
        roomName,
        messageOwner,
        messageID,
        content,
    ).Scan(&id, &creationDatetime, &editionDatetime); err == sql.ErrNoRows {
        return errors.Validation(
            errors.FromString("Apparently the message does'nt exists or aren't owned by you"))
    } else if err != nil {
        return errors.Internal(err)
    }
    go m.Broadcast(eventEdition, &message.Data{
        ID:               id,
        Content:          content,
        OwnerUsername:    messageOwner,
        RoomName:         roomName,
        RoomUsername:     username,
        CreationDatetime: creationDatetime,
        EditionDatetime:  editionDatetime,
    })
    return nil
}

func (m *model) Delete(username, roomName, messageID, messageOwner string) errors.Error {
    const query = `
    WITH (
        SELECT m.id as "msg_id", r.id as "room_id", u.id as "user_id"
                FROM message m
                INNER JOIN room r ON r.id = m.room_id
                INNER JOIN "user" u ON u.id = m.user_id
                WHERE
                    m.id = $4
                    u.username = $3
                    AND r.name = $1 AND r.username = $2
    ) AS "t"
    DELETE FROM "message"
        WHERE "message".id = t.msg_id
        AND message.room_id = t.room_id
        AND message.user_id = t.user_id
    RETURNING id;`

    var id string
    if err := m.db.QueryRow(query, username, roomName, messageOwner, messageID).Scan(&id); err == sql.ErrNoRows {
        return errors.Empty(err)
    } else if err != nil {
        return errors.Internal(err)
    }
    go m.Broadcast(eventDeletion, &message.Data{
        ID: id,
        RoomUsername:username,
        RoomName:roomName,
        OwnerUsername:messageOwner,

    })
    return nil
}

const (
    eventDeletion = "deletion"
    eventEdition  = "edition"
    eventCreation = "creation"
)

func (m *model) Broadcast(event string, data *message.Data) {
    key := mapKey{username: data.RoomUsername, roomName: data.RoomName}
    eventData := &message.EventData{Data: *data, Event: event}
    m.RLock()
    defer m.RUnlock()
    listeners := m.listeners[key]
    for _, ch := range listeners {
        ch <- eventData
    }
}

func (m *model) Listen(username, roomName string) <-chan *message.EventData {
    key := mapKey{username, roomName}
    ch := make(chan *message.EventData, 10)
    m.Lock()
    defer m.Unlock()
    list, ok := m.listeners[key]
    if !ok {
        list = make([](chan *message.EventData), 0, 31)
    }
    m.listeners[key] = append(list, ch)
    return ch
}

func (m *model) StopListening(username, roomName string, ch <-chan *message.EventData) {
    key := mapKey{username, roomName}
    m.Lock()
    defer m.Unlock()
    list := m.listeners[key]
    if list == nil {
        return
    }
    for k, v := range list {
        if v == ch {
            list = append(list[:k], list[k+1:]...)
            m.listeners[key] = list
            close(v)
            break
        }
    }
}

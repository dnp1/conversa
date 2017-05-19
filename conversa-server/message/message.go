package message

import (
    "time"
    "database/sql"
    "strings"
    "errors"
)

var (//errors
    ErrMessageIsEmpty = errors.New("Message is empty, invalid!")
    ErrCouldNotFindRoom = errors.New("Trying to create message on a unexisting room!")
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
    Edit(username, roomName, messageOwner, message, content string) error
    Delete(username, roomName, message, messageOwner string) error
    All(username , roomName string, limit, offset int64) ([]Data, error)
}

type Builder struct {
    DB *sql.DB
}


func (builder Builder) Build() Message {
    return &message{
        db: builder.DB,
    }
}

type message struct {
    db *sql.DB
}

func (m *message) Create(username , roomName, senderName, content string) error {
    const query = `INSERT INTO "message"("room_id", "user_id", "content")
    SELECT r.id, u.id, $4
        FROM room r
        LEFT JOIN "user" u ON u.username = $3
        WHERE r.name = $1 AND r.username = $2
    RETURNING id;`

    var id int64
    content = strings.TrimSpace(content)
    if content == "" {
        return ErrMessageIsEmpty
    }

    if err := m.db.QueryRow(query, username, roomName, senderName, content).Scan(&id); err == sql.ErrNoRows {
        return ErrCouldNotFindRoom
    } else if err != nil {
        return err
    }

    return nil
}

func (m *message) Edit(username, roomName, messageOwner, message, content string) error{
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
    UPDATE "message" SET "content" = $5
        WHERE "message".id = t.msg_id
        AND message.room_id = t.room_id
        AND message.user_id = t.user_id
    RETURNING id;`

    var id int64
    content = strings.TrimSpace(content)
    if content == "" {
        return ErrMessageIsEmpty
    }

    if err := m.db.QueryRow(query, username, roomName, messageOwner, message, content).Scan(&id); err == sql.ErrNoRows {
        return ErrCouldNotFindRoom
    } else if err != nil {
        return err
    }

    return nil
}

func (m *message) Delete(username, roomName, message, messageOwner string) error{
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

    var id int64
    if err := m.db.QueryRow(query, username, roomName, messageOwner, message).Scan(&id); err == sql.ErrNoRows {
        return ErrCouldNotFindRoom
    } else if err != nil {
        return err
    }

    return nil
}

func (m *message) All(username , roomName string, limit, offset int64) ([]Data, error) {

    return nil,nil
}


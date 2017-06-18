package message

import (
    "database/sql"
    "strings"
    "github.com/dnp1/conversa/server/errors"
    "github.com/dnp1/conversa/server/data/message"
)

var (//errors
    ErrMessageIsEmpty = errors.Validation(errors.FromString("Message is empty, invalid!"))
)


func New(db *sql.DB) *model {
    return &model{
        db: db,
    }
}

type model struct {
    db *sql.DB
}

func (m *model) Create(username , roomName, senderName, content string) errors.Error {
    const query = `INSERT INTO "model"("room_id", "user_id", "content")
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
        return errors.Validation(err)
    } else if err != nil {
        return errors.Internal(err)
    }

    return nil
}

func (m *model) Edit(username, roomName, messageOwner, message, content string) errors.Error {
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
    UPDATE "model" SET "content" = $5
        WHERE "model".id = t.msg_id
        AND model.room_id = t.room_id
        AND model.user_id = t.user_id
    RETURNING id;`

    var id int64
    content = strings.TrimSpace(content)
    if content == "" {
        return ErrMessageIsEmpty
    }

    if err := m.db.QueryRow(query, username, roomName, messageOwner, message, content).Scan(&id); err == sql.ErrNoRows {
        return errors.Validation(
            errors.FromString("Apparently the message does'nt exists or aren't owned by you"))
    } else if err != nil {
        return errors.Internal(err)
    }

    return nil
}

func (m *model) Delete(username, roomName, message, messageOwner string) errors.Error {
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
    DELETE FROM "model"
        WHERE "model".id = t.msg_id
        AND model.room_id = t.room_id
        AND model.user_id = t.user_id
    RETURNING id;`

    var id int64
    if err := m.db.QueryRow(query, username, roomName, messageOwner, message).Scan(&id); err == sql.ErrNoRows {
        return errors.Empty(err)
    } else if err != nil {
        return errors.Internal(err)
    }

    return nil
}

func (m *model) All(username , roomName string, limit, offset int64) ([]message.Data, errors.Error) {

    return nil,nil
}


package session

import (
    "database/sql"
    "golang.org/x/crypto/bcrypt"

    "github.com/twinj/uuid"
    "github.com/dnp1/conversa/server/errors"
    "github.com/dnp1/conversa/server/data/session"
)


func  New(db *sql.DB) *model {
    return &model{
        db: db,
    }
}

type model struct {
    db *sql.DB
}

func (s *model) Create(username string, password string) (string, errors.Error) {
    var (
        hashedPassword string
        userID int64
    )
    const selQuery = `SELECT password, id FROM "user" WHERE username = $1;`
    if err := s.db.QueryRow(selQuery, username).Scan(&hashedPassword, &userID); err == sql.ErrNoRows {
        return "", errors.Empty(err)
    } else if err != nil {
        return "", errors.Internal(err)
    } else  if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
        return "", errors.Validation(err)
    }

    key := uuid.NewV4().String()
    const insQuery = `INSERT INTO "user_session"("session_key", "user_id") VALUES($1, $2);`
    if _, err := s.db.Exec(insQuery, key, userID); err != nil {
        return "", errors.Internal(err)
    }

    return key, nil
}

func (s *model) Delete(token string) errors.Error {
    const query = `DELETE FROM "user_session" WHERE session_key = $1 RETURNING "user_id"`
    var id int64
    if err := s.db.QueryRow(query, token).Scan(&id); err != nil {
        if err == sql.ErrNoRows {
            return errors.Empty(err)
        }
        return errors.Internal(err)
    }
    return nil
}

func (s *model) Retrieve(token string) (*session.Data, errors.Error) {
    var data session.Data

    const query = `SELECT u.username, s.user_id
        FROM "user_session" s
        INNER JOIN "user" u ON s."user_id" = u."id"
 WHERE session_key = $1
        `
    if err := s.db.QueryRow(query, token).Scan(&data.Username, &data.UserID); err != nil {
        if err == sql.ErrNoRows {
            return  nil, errors.Empty(err)
        } else {
            return nil, errors.Internal(err)
        }
    }
    return &data, nil
}


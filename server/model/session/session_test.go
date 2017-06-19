package session_test

import (
    "testing"
    "github.com/dnp1/conversa/server/model/session"
    sessionData "github.com/dnp1/conversa/server/data/session"
    "gopkg.in/DATA-DOG/go-sqlmock.v1"
    "database/sql"
    "github.com/stretchr/testify/assert"
    "golang.org/x/crypto/bcrypt"

    "github.com/twinj/uuid"
    "github.com/dnp1/conversa/server/errors"
)

func TestSession_Create(t *testing.T) {
    var err errors.Error
    db, mock, err0 := sqlmock.New()
    assert.NoError(t, err0)
    s := session.New(db)
    username, password := "user", "password"
    //case 0, can't found
    {
        mock.ExpectQuery(".*").WithArgs(username).WillReturnError(sql.ErrNoRows)
        _, err = s.Create(username, password)
        assert.True(t, err.Empty())
    }

    //case 1
    {
        mock.ExpectQuery(".*").WithArgs(username).WillReturnError(errors.FromString("unexpected error"))
        _, err = s.Create(username, password)
        assert.Error(t, err)
    }

    columns := []string{"password", "user_id"}
    //case 2
    {

        rows := sqlmock.NewRows(columns).AddRow("1", 1)
        mock.ExpectQuery(".*").WithArgs(username).WillReturnRows(rows)
        _, err = s.Create(username, password)
    }
    assert.Error(t, err)
    //case 3
    {
        hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        rows := sqlmock.NewRows(columns).AddRow(hashed, 1)
        mock.ExpectQuery(".*").WithArgs(username).WillReturnRows(rows)
        mock.ExpectExec(".*").WithArgs(sqlmock.AnyArg(), 1).WillReturnError(errors.FromString("random error"))
        _, err = s.Create(username, password)
        assert.Error(t, err)
    }

    //case 4
    {
        hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        rows := sqlmock.NewRows(columns).AddRow(hashed, 1)
        mock.ExpectQuery(".*").WithArgs(username).WillReturnRows(rows)
        mock.ExpectExec(".*").WithArgs(sqlmock.AnyArg(), 1).WillReturnResult(sqlmock.NewResult(1, 1))
        _, err = s.Create(username, password)
        assert.NoError(t, err)
        assert.NoError(t, mock.ExpectationsWereMet())
    }
}

func TestSession_Delete(t *testing.T) {
    var err error
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    s := session.New(db)

    //case 0
    {
        token := uuid.NewV4().String()
        mock.ExpectQuery(".*").WillReturnError(errors.FromString("unexpected error"))
        assert.Error(t, s.Delete(token))
    }
    //case 1
    {
        token := uuid.NewV4().String()
        rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
        mock.ExpectQuery(".*").WillReturnRows(rows)
        assert.NoError(t, s.Delete(token))
        assert.NoError(t, mock.ExpectationsWereMet())
    }
}

func TestSession_Retrieve(t *testing.T) {
    var data = new(sessionData.Data)
    var err error

    var token string
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    s := session.New(db)

    //case 0
    {
        token = uuid.NewV4().String()
        mock.ExpectQuery(".*").WillReturnError(errors.FromString("unexpected error"))
        data, err := s.Retrieve(token)
        assert.Nil(t, data)
        assert.Error(t, err)
        assert.True(t, err.Internal())
    }

    //case 1
    {
        token = uuid.NewV4().String()
        mock.ExpectQuery(".*").WillReturnError(sql.ErrNoRows)
        data, err := s.Retrieve(token)
        assert.Nil(t, data)
        assert.Error(t, err)
        assert.True(t, err.Empty())
    }

    //case 2
    {
        token = uuid.NewV4().String()

        rows := sqlmock.NewRows([]string{"username", "user_id"}).AddRow("121212", 666)
        mock.ExpectQuery(".*").WillReturnRows(rows)
        data, err = s.Retrieve(token)
        assert.Nil(t, err)
        assert.NotNil(t, data)
        assert.NoError(t, mock.ExpectationsWereMet())
    }
}

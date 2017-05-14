package session_test

import (
    "testing"
    "github.com/dnp1/conversa/server/session"
    "gopkg.in/DATA-DOG/go-sqlmock.v1"
    "database/sql"
    "github.com/stretchr/testify/assert"
    "errors"
    "golang.org/x/crypto/bcrypt"
)

func TestSession_Create(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    s := session.Builder{DB:db}.Build()
    username, password := "user", "password"
    //case 0
    mock.ExpectQuery(".*").WithArgs(username).WillReturnError(sql.ErrNoRows)
    _, err = s.Create(username, password)
    assert.Equal(t, session.ErrBadCredentials,err )
    //case 1
    mock.ExpectQuery(".*").WithArgs(username).WillReturnError(errors.New("unexpected error"))
    _, err = s.Create(username, password)
    assert.Error(t, err)
    //case 2
    columns := []string{"password", "user_id"}
    rows := sqlmock.NewRows(columns).AddRow("1", 1)
    mock.ExpectQuery(".*").WithArgs(username).WillReturnRows(rows)
    _, err = s.Create(username, password)
    assert.Error(t, err)
    //case 3
    hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    rows = sqlmock.NewRows(columns).AddRow(hashed, 1)
    mock.ExpectQuery(".*").WithArgs(username).WillReturnRows(rows)
    mock.ExpectExec(".*").WithArgs(sqlmock.AnyArg(), 1).WillReturnError(errors.New("random error"))
    _, err = s.Create(username, password)
    assert.Error(t, err)
    //case 4
    hashed, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    rows = sqlmock.NewRows(columns).AddRow(hashed, 1)
    mock.ExpectQuery(".*").WithArgs(username).WillReturnRows(rows)
    mock.ExpectExec(".*").WithArgs(sqlmock.AnyArg(), 1).WillReturnResult(sqlmock.NewResult(1,1))
    _, err = s.Create(username, password)
    assert.NoError(t, err)
}
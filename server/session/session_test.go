package session_test

import (
    "testing"
    "github.com/dnp1/conversa/server/session"
    "gopkg.in/DATA-DOG/go-sqlmock.v1"
    "database/sql"
    "github.com/stretchr/testify/assert"
    "errors"
    "golang.org/x/crypto/bcrypt"
    "github.com/satori/go.uuid"
)

func TestSession_Create(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    s := session.Builder{DB:db}.Build()
    username, password := "user", "password"
    //case 0
    mock.ExpectQuery(".*").WithArgs(username).WillReturnError(sql.ErrNoRows)
    _, err = s.Create(username, password)
    assert.Equal(t, session.ErrBadCredentials, err)
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
    mock.ExpectExec(".*").WithArgs(sqlmock.AnyArg(), 1).WillReturnResult(sqlmock.NewResult(1, 1))
    _, err = s.Create(username, password)
    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSession_Delete(t *testing.T) {
    var err error
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    s := session.Builder{DB:db}.Build()

    //case 0
    token := uuid.NewV4().String()
    mock.ExpectExec(".*").WillReturnError(errors.New("unexpected error"))
    assert.Error(t, s.Delete(token))
    //case 1
    token = uuid.NewV4().String()
    mock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(1, 1))
    assert.NoError(t, s.Delete(token))
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSession_Valid(t *testing.T) {
    var token string
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    s := session.Builder{DB:db}.Build()

    //case 0
    token = uuid.NewV4().String()
    mock.ExpectQuery(".*").WillReturnError(errors.New("unexpected error"))
    assert.Error(t, s.Valid(token))

    //case 1
    token = uuid.NewV4().String()
    mock.ExpectQuery(".*").WillReturnError(sql.ErrNoRows)
    assert.EqualValues(t, session.ErrTokenNotFound, s.Valid(token))

    //case 2 (in theory impossible, but...)
    token = uuid.NewV4().String()
    columns := []string{"a"}
    rows := sqlmock.NewRows(columns)
    rows.AddRow(false)
    mock.ExpectQuery(".*").WillReturnRows(rows)
    assert.EqualValues(t, session.ErrTokenNotFound, s.Valid(token))

    //case 3
    token = uuid.NewV4().String()
    rows = sqlmock.NewRows(columns)
    rows.AddRow(true)
    mock.ExpectQuery(".*").WillReturnRows(rows)
    assert.NoError(t, s.Valid(token))
    assert.NoError(t, mock.ExpectationsWereMet())

}


func TestSession_Retrieve(t *testing.T) {
    var data  = new(session.Data)
    var err error

    var token string
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    s := session.Builder{DB:db}.Build()

    //case 0
    token = uuid.NewV4().String()
    mock.ExpectQuery(".*").WillReturnError(errors.New("unexpected error"))
    data, err = s.Retrieve(token)
    assert.Nil(t, data)
    assert.Error(t, err)

    //case 1
    token = uuid.NewV4().String()
    mock.ExpectQuery(".*").WillReturnError(sql.ErrNoRows)
    data, err = s.Retrieve(token)
    assert.Nil(t, data)
    assert.Equal(t, session.ErrTokenNotFound, err)

    //case 2
    token = uuid.NewV4().String()

    rows := sqlmock.NewRows([]string{"username", "user_id"}).AddRow("121212", 666)
    mock.ExpectQuery(".*").WillReturnRows(rows)
    data, err = s.Retrieve(token)
    assert.Nil(t, err)
    assert.NotNil(t, data)
    assert.NoError(t, mock.ExpectationsWereMet())
}
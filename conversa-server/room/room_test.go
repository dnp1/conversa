package room_test

import (
    "testing"
    "gopkg.in/DATA-DOG/go-sqlmock.v1"
    "github.com/stretchr/testify/assert"
    "github.com/dnp1/conversa/conversa-server/room"
    "github.com/pkg/errors"
    "database/sql"
)

func TestRoom_Create(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    r := room.Builder{DB:db}.Build()
    username, roomname := "user", "room"
    //case 0
    mock.ExpectQuery(".*").WillReturnError(errors.New("unexpected"))
    assert.Error(t, r.Create(username, roomname))
    //case 1
    mock.ExpectQuery(".*").WillReturnError(sql.ErrNoRows)
    assert.Error(t, r.Create(username, roomname))
    //case 2
    mock.ExpectQuery(".*").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
    assert.NoError(t, r.Create(username, roomname))
    assert.NoError(t, mock.ExpectationsWereMet())
}


func TestRoom_Delete(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    r := room.Builder{DB:db}.Build()
    username, roomname := "user", "room"
    //case 0
    mock.ExpectExec(".*").WillReturnError(errors.New("unexpected!"))
    assert.Error(t, r.Delete(username, roomname))
    //case 1
    mock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(1,1))
    assert.NoError(t, r.Delete(username, roomname))
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoom_Rename(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    r := room.Builder{DB:db}.Build()
    username, roomname := "user", "room"
    //case 0
    mock.ExpectExec(".*").WillReturnError(errors.New("unexpected!"))
    assert.Error(t, r.Rename(username, roomname, "new"))
    //case 1
    mock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(1,1))
    assert.NoError(t, r.Rename(username, roomname, "new"))
    assert.NoError(t, mock.ExpectationsWereMet())
}




func TestRoom_All(t *testing.T) {
    var data []room.Data
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    r := room.Builder{DB:db}.Build()
    //case 0
    mock.ExpectQuery(".*").WillReturnError(errors.New("unexpected!"))
    data, err = r.All()
    assert.Error(t, err)
    assert.Nil(t,data)
    //case 1
    rows := sqlmock.NewRows([]string{"username"})
    rows.AddRow("a")
    rows.AddRow("danilo")
    mock.ExpectQuery(".*").WillReturnRows(rows)
    data, err = r.All()
    assert.Error(t, err)
    assert.Nil(t,data)
    //case 2
    rows = sqlmock.NewRows([]string{"username","name"})
    rows.AddRow("a","b")
    rows.AddRow("danilo","programação")
    mock.ExpectQuery(".*").WillReturnRows(rows)
    data, err = r.All()
    assert.NoError(t, err)
    assert.NotNil(t,data)
    assert.NoError(t, mock.ExpectationsWereMet())
}


func TestRoom_AllByUser(t *testing.T) {
    var data []room.Data
    var user = "fulano"
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    r := room.Builder{DB:db}.Build()
    //case 0
    mock.ExpectQuery(".*").WillReturnError(errors.New("unexpected!"))
    data, err = r.AllByUser(user)
    assert.Error(t, err)
    assert.Nil(t,data)
    //case 1
    rows := sqlmock.NewRows([]string{"username", "b"})
    rows.AddRow("a", "a")
    rows.AddRow("danilo", "dsada")
    mock.ExpectQuery(".*").WillReturnRows(rows)
    data, err = r.AllByUser(user)
    assert.Error(t, err)
    assert.Nil(t,data)
    //case 2
    rows = sqlmock.NewRows([]string{"name"})
    rows.AddRow("a")
    rows.AddRow("danilo")
    mock.ExpectQuery(".*").WillReturnRows(rows)
    data, err = r.AllByUser(user)
    assert.NoError(t, err)
    assert.NotNil(t,data)
    assert.NoError(t, mock.ExpectationsWereMet())
}


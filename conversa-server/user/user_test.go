package user_test

import (
    "testing"
    "gopkg.in/DATA-DOG/go-sqlmock.v1"
    "github.com/dnp1/conversa/conversa-server/user"
    "github.com/stretchr/testify/assert"
    "github.com/pkg/errors"
    "golang.org/x/crypto/bcrypt"

    "github.com/twinj/uuid"
    "strings"
)

func TestUser_Create(t *testing.T) {
    var cost int = bcrypt.DefaultCost
    var pass = uuid.NewV4().String()
    //case 0 bcrypt error
    assert.Error(t, user.Builder{BCryptCost: 2 << 15}.Build().Create("user", pass, pass))
    //diferent password and confirmation
    u := user.Builder{BCryptCost:cost}.Build()
    assert.Equal(t, user.ErrPasswordConfirmationDoesNotMatch, u.Create("user", pass, uuid.NewV4().String()))
    //case sql error
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    u = user.Builder{DB: db, BCryptCost:cost}.Build()
    expectedErr := errors.New("puts")
    mock.ExpectQuery(".*").WillReturnError(expectedErr).WithArgs("user", sqlmock.AnyArg())
    assert.Equal(t, expectedErr, u.Create("user", pass, pass))
    //taken
    mock.ExpectQuery(".*").WillReturnRows(sqlmock.NewRows([]string{"id"})).WithArgs("user", sqlmock.AnyArg())
    assert.Equal(t, user.ErrUsernameAlreadyTaken, u.Create("user", pass, pass))
    //ok
    mock.ExpectQuery(".*").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1)).WithArgs("user", sqlmock.AnyArg())
    assert.NoError(t, u.Create("user", pass, pass))
    assert.NoError(t, mock.ExpectationsWereMet())

}


func TestValidate(t *testing.T) {
    type Input struct {
        Username string
        Password string
        PasswordConfirmation string
    }
    type Case struct {
        Input Input
        Output error
    }
    validName := strings.Repeat("k", (user.UsernameMaxLength+user.UsernameMaxLength)/2)
    cases := []Case{
        {
            Input{validName,"aang", "ang"},
            user.ErrPasswordConfirmationDoesNotMatch,
        },
        {
            Input{strings.Repeat("k", user.UsernameMaxLength+1),"aang", "aang"},
            user.ErrUsernameWrongLength,
        },
        {
            Input{strings.Repeat("c", user.UsernameMinLength-1),"aang", "aang"},
            user.ErrUsernameWrongLength,
        },
        {
            Input{validName, "aang", "aang"},
            user.ErrPasswordTooShort,
        },
        {
            Input{strings.Repeat("c", user.UsernameMinLength)+"*", "wdasasdasdasdada", "wdasasdasdasdada"},
            user.ErrUsernameHasInvalidCharacters,
        },
    }

    for i, c := range cases {
        out := user.Validate(
            c.Input.Username,
            c.Input.Password,
            c.Input.PasswordConfirmation,
        )
        if !assert.Equal(t, c.Output, out) {
            t.Logf("Case %d faild!%+v\n", i, c.Input)
        }
    }
}
package server_test

import (
    "testing"

    "net/http/httptest"
    "net/http"
    "github.com/stretchr/testify/assert"
    "github.com/dnp1/conversa/server"
    "strings"
    "io"
    "fmt"
    "gopkg.in/gin-gonic/gin.v1"
)

func init() {
    gin.SetMode(gin.TestMode)
}

//TestLogin0 cases that end in BadRequest
func TestLogin0(t *testing.T) {

    router := server.NewRouter()
    cases := [...]io.Reader {
        nil,
        strings.NewReader(`{"user_name": "json","password"}`),
        strings.NewReader(fmt.Sprintf(
            `{"username": "%s", "password": "passphrase"}`,
            strings.Repeat("110111001011101111000100110101011110011011110", 10),
        )),
        strings.NewReader(fmt.Sprintf(
            `{"username": "user", "password": "%s"}`,
            strings.Repeat("110111001011abc101111000100110101011110011011110", 10),
        )),
        strings.NewReader(`{"username": "user", "password": "passphrase"}`),
    }
    for i, c := range cases {
        req, err := http.NewRequest("POST", "/session", c)
        if assert.NoError(t, err) {
            resp :=  httptest.NewRecorder()
            router.ServeHTTP(resp, req)
            if !assert.Exactly(t, http.StatusBadRequest, resp.Code){
                t.Logf("Case %s", i)
            }
        }
    }
}

//TestLogin0 cases that end in Ok
func TestLogin1(t *testing.T) {
    router := server.NewRouter()
    body := strings.NewReader(`{"username":"user", "passwor":"passphrase"`)
    req, err := http.NewRequest("POST", "/session", body)
    if assert.NoError(t, err) {
        resp :=  httptest.NewRecorder()
        router.ServeHTTP(resp, req)
        assert.Exactly(t, http.StatusOK, resp.Code)
    }
}

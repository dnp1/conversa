package client

import (
    "time"
    "net/http"
    "errors"
    "io"
)

var transport = &http.Transport{
    TLSHandshakeTimeout:10 * time.Second,
}

var (
    //errors
    ErrInvalidTarget = errors.New("Invalid target!")
)

type RoomItem struct {
    Username string `json:"username"`
    Name     string `json:"name"`
}

type RoomData struct {
    Message string `json:"message"`
    Items []RoomItem `json:"data"`
}

type RoomBody struct {
    Name string `json:"name"`
}

type LoginBody struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type Requester interface {
    NewRequest(method, endpoint string, body io.Reader) (*http.Request)
    Do(req *http.Request, jar http.CookieJar) (*http.Response, error)
    Request(method, endpoint string, body io.Reader, jar http.CookieJar) (*http.Response, error)
}
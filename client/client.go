package client

import (
    "net/http"
    "net/url"
)

type Client interface {
    Login(username, password string) Error
    Logout() error
    RoomCreate(name string) error
    RoomList() ([]RoomData, error)
    RoomRemove(name string) error
    RoomRename(currentName string, newName string) error
    JSON() string
    //SignUp(username, password, passwordConfirmation string) Error
}

type ClientBuilder struct {
    Target string `json:"target"`
    Username string
    Cookies []*http.Cookie
}

func (builder ClientBuilder) Build() (Client, Error) {
    if requester, err := newRequester(builder.Target, builder.Cookies); err != nil{
        return nil, err
    } else {
        return &client{
            requester: requester,
        }, nil
    }
}



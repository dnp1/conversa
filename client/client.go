package client

import (
    "net/http"
    "fmt"
    "github.com/dnp1/conversa/client/requester"
    "github.com/dnp1/conversa/client/errors"
)

type Credentials struct {
    Username string
    Cookies []*http.Cookie
}

func New(target string, credentials []byte) (*cl, requester.Error) {
    //TODO:Treat credentials
    if req, err := requester.New(target, nil); err != nil{
        return nil, err
    } else {
        return &cl{
            requester: req,
        }, nil
    }
}

var ErrNotLoggedIn = errors.Authentication(
    errors.FromString("To perform this action you must provide credentials"))

type Requester interface {
    Request(method, path string, body, resp interface {}) requester.Error
}

type cl struct {
    username string
    requester Requester
}

func (cl *cl) Login(username, password string) requester.Error {
    const path = "/sessions"
    body := LoginBody{Username:username, Password:password}
    if err := cl.requester.Request(http.MethodPost, path, body, &EmptyResponse{}); err != nil {
        return err
    }
    cl.username = username
    return nil
}

func (cl *cl) Logout() requester.Error {
    const path = "/session"
    if err := cl.requester.Request(http.MethodDelete, path, nil, &EmptyResponse{}); err != nil {
        return err
    }
    cl.username = ""
    return nil
}

func (cl *cl) RoomCreate(name string) requester.Error {
    if cl.username == "" {
        return ErrNotLoggedIn
    }
    path := fmt.Sprintf("/users/%s/room", cl.username)
    body := RoomBody{Name:name}
    return cl.requester.Request(http.MethodPost, path, body, &EmptyResponse{})
}

func (cl *cl) RoomList() ([]RoomItem,  requester.Error) {
    if cl.username == "" {
        return nil, ErrNotLoggedIn
    }
    const path = "/room"
    var response RoomData
    if err := cl.requester.Request(http.MethodGet, path, nil, &response); err != nil {
        return nil, err
    }
    return response.Items, nil
}

func (cl *cl) RoomRemove(name string) requester.Error {
    if cl.username == "" {
        return ErrNotLoggedIn
    }
    path := fmt.Sprintf("/users/%s/room/%s", cl.username, name)
    return cl.requester.Request(http.MethodDelete, path, nil, &EmptyResponse{})
}

func (cl *cl) RoomRename(currentName string, newName string) requester.Error {
    if cl.username == "" {
        return ErrNotLoggedIn
    }
    path := fmt.Sprintf("/users/%s/room/%s", cl.username, currentName)
    body := RoomBody{Name:newName}
    return cl.requester.Request(http.MethodPatch, path, body, &EmptyResponse{})
}

func (cl *cl) SignUp(username, password, passwordConfirmation string) requester.Error {
    path := fmt.Sprintf("/user")
    body := SignUp{
        Username:username,
        Password:password,
        PasswordConfirmation: passwordConfirmation,
    }
    return cl.requester.Request(http.MethodPost, path, body, &EmptyResponse{})
}
package client

import (
    "net/http"
    "fmt"
    "github.com/dnp1/conversa/client/requester"
    "github.com/dnp1/conversa/client/errors"
    "encoding/json"
    "io/ioutil"
)

type credentials struct {
    Username string `json:"username"`
    Cookies  []*http.Cookie `json:"cookies"`
}

func New(target string, cred []byte) (*cl, errors.Error) {
    var data credentials
    var username string
    if cred != nil {
        if err := json.Unmarshal(cred, &data); err != nil {
            const msg = "Invalid credentials provided. Maybe you should login again or delete old credentials"
            return nil, errors.Unexpected(errors.FromString(msg))
        }
        username = data.Username
    }
    if req, err := requester.New(target, data.Cookies); err != nil {
        return nil, err
    } else {
        return &cl{
            requester: req,
            username:  username,
        }, nil
    }
}

var ErrNotLoggedIn = errors.Authentication(
    errors.FromString("To perform this action you must provide credentials"))


type Requester interface {
    Request(method, path string, body, resp interface{}) errors.Error
    Cookies() []*http.Cookie
    NotifySSE(path string, evCh chan<- *requester.Sse, errCh chan <- errors.Error)
}

type cl struct {
    username  string
    requester Requester
}


func (cl *cl) Credentials() []byte {
    bytes, _ :=  json.Marshal(
        credentials{
            Username: cl.username,
            Cookies:  cl.requester.Cookies(),
        })
    return bytes
}

func (cl *cl) Login(username, password string) errors.Error {
    const path = "/session"
    body := LoginBody{Username: username, Password: password}
    if err := cl.requester.Request(http.MethodPost, path, body, &EmptyResponse{}); err != nil {
        return err
    }
    cl.username = username
    return nil
}

func (cl *cl) Logout() errors.Error {
    const path = "/session"
    if err := cl.requester.Request(http.MethodDelete, path, nil, &EmptyResponse{}); err != nil {
        return err
    }
    cl.username = ""
    return nil
}

func (cl *cl) RoomCreate(name string) errors.Error {
    if cl.username == "" {
        return ErrNotLoggedIn
    }
    path := fmt.Sprintf("/user/%s/room", cl.username)
    body := RoomBody{Name: name}
    return cl.requester.Request(http.MethodPost, path, body, &EmptyResponse{})
}

func (cl *cl) RoomList() ([]RoomItem, errors.Error) {
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

func (cl *cl) RoomRemove(name string) errors.Error {
    if cl.username == "" {
        return ErrNotLoggedIn
    }
    path := fmt.Sprintf("/user/%s/room/%s", cl.username, name)
    return cl.requester.Request(http.MethodDelete, path, nil, &EmptyResponse{})
}

func (cl *cl) SignUp(username, password, passwordConfirmation string) errors.Error {
    path := fmt.Sprintf("/user")
    body := SignUp{
        Username:             username,
        Password:             password,
        PasswordConfirmation: passwordConfirmation,
    }
    return cl.requester.Request(http.MethodPost, path, body, &EmptyResponse{})
}

func (cl *cl) MessageCreate(user, room, content string) errors.Error {
    if cl.username == "" {
        return ErrNotLoggedIn
    }
    path := fmt.Sprintf("/user/%s/room/%s/message", user, room)
    body := MessageBody{
        Content: content,
    }
    return cl.requester.Request(http.MethodPost, path, body, &EmptyResponse{})
}

func (cl *cl) MessageEdit(user, room, messageId, content string) errors.Error {
    if cl.username == "" {
        return ErrNotLoggedIn
    }
    path := fmt.Sprintf("/user/%s/room/%s/message/%s", user, room, messageId)
    body := MessageBody{
        Content: content,
    }
    return cl.requester.Request(http.MethodPatch, path, body, &EmptyResponse{})
}

func (cl *cl) MessageDelete(user, room, messageId string) errors.Error {
    if cl.username == "" {
        return ErrNotLoggedIn
    }
    path := fmt.Sprintf("/user/%s/room/%s/message/%s", user, room, messageId)
    return cl.requester.Request(http.MethodDelete, path, nil, &EmptyResponse{})
}

func (cl *cl) Listen(user, room string, ch chan <- *Message, errCh chan <- errors.Error) {
    if cl.username == "" {
        errCh <- ErrNotLoggedIn
    }
    path := fmt.Sprintf("/user/%s/room/%s/listen", user, room)
    sseCh := make(chan *requester.Sse)
    sseErrCh := make(chan errors.Error)
    go cl.requester.NotifySSE(path, sseCh, sseErrCh)
    for {
        select {
        case err, ok := <-sseErrCh:
            if ok {
                close(sseCh)
                close(sseErrCh)
                errCh <- err
            }
            break
        case sse, ok := <-sseCh:
            if ok {
                var data = new(Message)
                if bytes, err := ioutil.ReadAll(sse.Data); err != nil {
                    errCh <- errors.Unexpected(err)
                    break
                } else if err := json.Unmarshal(bytes, data); err != nil {
                    errCh <- errors.Unexpected(err)
                    break
                }
                ch <- data
            } else {
                break
            }
        }
    }

}



package client

import (
    "net/http/cookiejar"
    "net/http"
    "net/url"
    "log"
    "encoding/json"
    "fmt"
    "bytes"
    "io/ioutil"
    "io"
    "github.com/dnp1/conversa/conversa-cli/cl"
)

type Session interface {
    Logout() error
    RoomCreate(name string) error
    RoomList() ([]RoomData, error)
    RoomRemove(name string) error
    RoomRename(currentName string, newName string) error
}

type SessionBuilder struct {
    Requester Requester
    Target    string
    Username  string
    Cookies   []*http.Cookie
}

func (builder *SessionBuilder) Build() Session {
    if jar, err := cookiejar.New(nil); err != nil {
        log.Fatalln(err)
    } else if u, err := url.Parse(builder.Target); err != nil {
        log.Fatalln(err)
    } else {
        jar.SetCookies(u, builder.Cookies)
        return &session{
            requester:builder.Requester,
            username:builder.Username,
            jar: jar,
        }
    }
    return Session(nil)

}

type session struct {
    requester Requester
    jar       http.CookieJar
    username  string
}

func (s *session) Logout() error {

}

func (s *session) RoomCreate(name string) error {
    body := RoomBody{Name:name}
    if js, err := json.Marshal(body); err != nil {
        return err //barely impossible
    } else {
        reader := bytes.NewReader(js)
        endpoint := fmt.Sprintf("users/%s/rooms", s.username)
        if resp, err := s.requester.Request(
            http.MethodPost,
            endpoint,
            reader,
            s.jar,
        ); err != nil {
            return err
        } else {
            defer resp.Body.Close()
            if resp.StatusCode == http.StatusCreated {
                return nil
            } else {
                //TODO:improve error handling
                return nil
            }
        }
    }
}

func (s *session) RoomList() ([]RoomItem, error) {
    const endpoint = "/rooms"
    if resp, err := s.requester.Request(
        http.MethodGet,
        endpoint,
        nil,
        s.jar,
    ); err != nil {
        return err
    } else {
        defer resp.Body.Close()
        if resp.StatusCode == http.StatusOK {
            var data RoomData
            if err := ReadJSON(resp.Body, &data); err != nil {
                return nil, err
            }
            return data.Items, nil
        } else {
            //TODO:improve error handling
            return nil, nil
        }
    }
}

func (s *session) RoomRemove(name string) error {
    endpoint := fmt.Sprintf("/users/%s/rooms/%s", s.username, name)
    if resp, err := s.requester.Request(
        http.MethodDelete,
        endpoint,
        nil,
        s.jar,
    ); err != nil {
        return err
    } else {
        defer resp.Body.Close()
        if resp.StatusCode != http.StatusOK {
            //TODO:improve error handling
            return fmt.Errorf("Error Status: %d", resp.StatusCode)
        }
    }
    return nil
}
func (s *session) RoomRename(currentName string, newName string) error {
    body := RoomBody{Name:newName}
    if js, err := json.Marshal(body); err != nil {
        return err //barely impossible
    } else if {
        endpoint := fmt.Sprintf("/users/%s/rooms/%s", s.username, currentName)
    }
}

func ReadJSON(body io.Reader, refToData interface{}) error {
    if data, err := ioutil.ReadAll(body); err != nil {
        return err
    } else if err := json.Unmarshal(data, refToData); err != nil {
        return err
    }
    return nil
}
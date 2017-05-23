package client

import (
    "net/http"
    "fmt"
    "errors"
)

type client struct {
    username string
    requester Requester
}

func (c *client) Login(username, password string) Error {
    body := LoginBody{Username:username, Password:password}
    if jsReader, err := JSONReader(body); err != nil {
        return err
    } else {
        const endpoint = "/sessions"
        if code, body, err :=
            c.requester.SimpleRequest(
                http.MethodPost,
                endpoint,
                jsReader,
            ); err != nil {
            return err
        } else {
            return HttpError(body, code)
        }
    }
}




func (s *client) Logout() error {
    return nil
}

func (s *client) RoomCreate(name string) error {
    body := RoomBody{Name:name}
    if reader, err := JSONReader(body); err != nil {
        return err //barely impossible
    } else {
        endpoint := fmt.Sprintf("users/%s/rooms", s.username)
        if resp, err := s.requester.Request(
            http.MethodPost,
            endpoint,
            reader,
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

func (s *client) RoomList() ([]RoomItem, error) {
    const endpoint = "/rooms"
    if resp, err := s.requester.Request(
        http.MethodGet,
        endpoint,
        nil,
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

func (s *client) RoomRemove(name string) error {
    endpoint := fmt.Sprintf("/users/%s/rooms/%s", s.username, name)
    if resp, err := s.requester.Request(
        http.MethodDelete,
        endpoint,
        nil,
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

func (s *client) RoomRename(currentName string, newName string) error {
    body := RoomBody{Name:newName}
    if reader, err := JSONReader(body); err != nil {
        return err
    } else if {
        endpoint := fmt.Sprintf("/users/%s/rooms/%s", s.username, currentName)
        if code, resp,  err := s.requester.SimpleRequest(
            http.MethodDelete,
            endpoint,
            reader,
        ); err != nil {
            return err
        } else if IsServerErrorCode(code) {
            return newServerError(errors.New(resp.Message))
        } else {

        }
    }
    return nil
}


package client

import (
    "net/http"
    "net/url"
    "errors"
)

type Client interface {
    Login(username, password string) (Session, error)
    SignUp(username, password, passwordConfirmation string) error
}

type Builder struct {
    Target string `json:"target"`
}

func (builder Builder) Build() (Client, Error) {
    if _, err := url.ParseRequestURI(builder.Target); err != nil {
        return nil, err
    }
    return &client{
        requester{target:builder.Target},
    }
}

type client struct {
    requester
}

func (c *client) Login(username, password string) (Session, Error) {
    body := LoginBody{Username:username, Password:password}
    if jsReader, err := JSONReader(body); err != nil {
        return err
    } else {
        const endpoint = "/sessions"
        if resp, err :=
            c.Request(
                http.MethodPost,
                endpoint,
                jsReader,
                nil,
            ); err != nil {
            return nil, err
        } else {
            if respBody, err := ReadResponseBody(resp.Body); err != nil {
                return nil, err
            } else {
                if code := resp.StatusCode; code == http.StatusCreated {
                    //TODO:Create_Session
                    return nil, nil
                } else if IsServerErrorCode(code) {
                    return newServer(errors.New(respBody.Message))
                }
            }
        }
    }
    return nil, nil
}

func IsServerErrorCode(code int) bool {
    return code >= 500 && code < 600
}




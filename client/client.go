package client

import (
    "encoding/json"
    "fmt"
    "bytes"
    "net/http"
    "time"
    "log"
    "io"
)

type Client interface {
    Login(username, password string) (Session, error)
    SignUp(username, password, passwordConfirmation string) error
}

type Builder struct {
    Target string `json:"target"`
}

func (builder Builder) Build() Client {
    return &client{
        target:builder.Target,
    }
}

type client struct {
    target string
}

func (c *client) NewRequest(method, endpoint string, body io.Reader) (*http.Request) {
    path := c.target + endpoint
    if req, err := http.NewRequest(method, path, body); err != nil {
        log.Fatalln(err)
        return nil
    } else {
        if body != nil {
            req.Header.Set("Content-Type", "application/json")
        }
        return req
    }
}

func (c *client) Do(req *http.Request, jar http.CookieJar) (*http.Response, error) {
    c := &http.Client{
        Transport:transport,
        Timeout: time.Second * 15,
        Jar: jar,
    }
    return c.Do(req)
}

func (c *client) Request(method, endpoint string, body io.Reader, jar http.CookieJar) (*http.Response, error) {
    req := c.NewRequest(method, endpoint, body)
    return c.Do(req, jar)
}

func (c *client) Login(username, password string) (Session, error) {
    if c.target == "" {
        return ErrInvalidTarget
    }
    body := LoginBody{Username:username, Password:password}
    if js, err := json.Marshal(body); err != nil {
        return err //barely impossible
    } else {
        bodyReader := bytes.NewReader(js)
        req := c.NewRequest(http.MethodPost, "/sessions", bodyReader, )
        if resp, err := c.Do(req, nil); err != nil {
            return nil, err
        } else {
            defer resp.Body.Close()
            if resp.StatusCode == http.StatusCreated {
                //TODO:Create_Session
                return nil, nil
            } else {
                return fmt.Errorf("Error Status: %d", resp.StatusCode)
            }
        }
    }
}





package cl

import (
    "net/http"
    "os"
    "io/ioutil"
    "encoding/json"
    "strings"
)

type Configuration struct {
    Path     string `json:"-"`
    Target   string `json:"target"`
    Username string `json:"username"`
    Cookies  []*http.Cookie `json:"cookie"`
}

func LoadFromPath(configPath string) (*Configuration, error) {
    var conf  = Configuration{Path: configPath}
    if _, err := os.Stat(configPath); err != nil {
        return nil, err
    } else if f, err := ioutil.ReadFile(configPath); err != nil {
        return nil, err
    } else if err := json.Unmarshal(f, &conf); err != nil {
        return nil, err
    } else {
        return &conf, nil
    }
}

func (c *Configuration) Create() error {
    if _, err := os.Stat(c.Path); err == nil {
        if err := os.Remove(c.Path); err != nil {
            return err
        }
    }
    if f, err := os.Create(c.Path); err != nil {
        return err
    } else if js, err := json.Marshal(c); err != nil {
        return err
    } else if _, err := f.Write(js); err != nil {
        return err
    } else {
        return f.Close()
    }
}

func (c *Configuration) SetTarget(target string) error {
    c.Target = strings.TrimSpace(target)
    return c.Create()
}

func (c *Configuration) SetSession(username string, cookies []*http.Cookie) error {
    c.Username = username
    c.Cookies = cookies
    return c.Create()
}

func (c*Configuration) UnsetSession() error {
    c.Username = ""
    c.Cookies = nil
    return c.Create()
}
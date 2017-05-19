package cl

import (
    "net/http"
    "time"
    "encoding/json"
    "fmt"
    "bytes"
    "net/http/cookiejar"
    "github.com/pkg/errors"
    "net/url"
    "io/ioutil"
)

const JsonContentType = "application/json"

type RoomData struct {
    Username string `json:"username"`
    Name     string `json:"name"`
}

type CL interface {
    Login(username, password string) error
    Logout() error
    RoomCreate(name string) error
    RoomList() ([]RoomData, error)
    RoomsUserList(username string) ([]RoomData, error)
    RoomRemove(name string) error
    RoomRename(currentName string, newName string) error
    SignUp(username, password, passwordConfirmation string) error
}

type client struct {
    config    *Configuration
    jar       http.CookieJar
    transport *http.Transport
}

type LoginBody struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func (builder *Configuration) BuildClient() (CL, error) {
    var cl = new(client)
    cl.config = builder
    cl.transport = &http.Transport{
        TLSHandshakeTimeout:10 * time.Second,
    }

    if builder.Cookies != nil {
        if jar, err := cookiejar.New(nil); err != nil {
            return nil, err
        } else if u, err := url.Parse(cl.config.Target); err != nil {
            return nil, err
        } else {
            jar.SetCookies(u, cl.config.Cookies)
            cl.jar = jar
        }
    }

    return cl, nil
}

var (
    ErrInvalidTarget = errors.New("Invalid target!")
)

func (cl *client) HttpClient() *http.Client {
    return &http.Client{
        Transport:cl.transport,
        Timeout: time.Second * 15,
        Jar: cl.jar,
    }
}

func (cl *client) Login(username, password string) error {
    if cl.config.Target == "" {
        return ErrInvalidTarget
    }
    if js, err := json.Marshal(LoginBody{Username:username, Password:password}); err != nil {
        return err //barely impossible
    } else {
        client := cl.HttpClient()
        url := fmt.Sprintf("%s/sessions", cl.config.Target)
        if resp, err := client.Post(url, JsonContentType, bytes.NewReader(js)); err != nil {
            return err
        } else {
            defer resp.Body.Close()
            if resp.StatusCode == http.StatusOK {
                cl.config.SetSession(username, resp.Cookies())
                return nil
            } else {
                return fmt.Errorf("Error Status: %d", resp.StatusCode)
            }
        }
    }
    return nil
}

func (cl *client) Logout() error {
    if cl.config.Target == "" {
        return ErrInvalidTarget
    }
    client := cl.HttpClient()
    url := fmt.Sprintf("%s/sessions", cl.config.Target)
    if req, err := http.NewRequest("DELETE", url, nil); err != nil {
        return err
    } else {

        if resp, err := client.Do(req); err != nil {
            return err
        } else {
            defer resp.Body.Close()
            if resp.StatusCode == http.StatusOK {
                return cl.config.UnsetSession()
            } else {
                return fmt.Errorf("Error Status: %d", resp.StatusCode)
            }
        }
    }
    return nil
}

type RoomBody struct {
    Name string `json:"name"`
}

func (cl *client) RoomCreate(name string) error {
    if cl.config.Target == "" {
        return ErrInvalidTarget
    }
    if js, err := json.Marshal(RoomBody{Name:name}); err != nil {
        return err //barely impossible
    } else {
        client := cl.HttpClient()
        url := fmt.Sprintf("%s/users/%s/rooms", cl.config.Target, cl.config.Username)
        if resp, err := client.Post(url, JsonContentType, bytes.NewReader(js)); err != nil {
            return err
        } else {
            defer resp.Body.Close()
            if resp.StatusCode == http.StatusOK {
                return nil
            } else {
                return fmt.Errorf("Error Status: %d", resp.StatusCode)
            }
        }
    }
    return nil
}

func (cl *client) RoomList() ([]RoomData, error) {
    if cl.config.Target == "" {
        return nil, ErrInvalidTarget
    }
    client := cl.HttpClient()
    url := fmt.Sprintf("%s/rooms", cl.config.Target)
    if resp, err := client.Get(url); err != nil {
        return nil, err
    } else {
        defer resp.Body.Close()
        if resp.StatusCode == http.StatusOK {
            var rooms []RoomData
            if data, err := ioutil.ReadAll(resp.Body); err != nil {
                return nil, err
            } else if err := json.Unmarshal(data, &rooms); err != nil {
                fmt.Println(string(data))
                return nil, err
            } else {
                return rooms, nil
            }
        } else {
            return nil, fmt.Errorf("Error Status: %d", resp.StatusCode)
        }
    }
}

func (cl *client) RoomsUserList(username string) ([]RoomData, error) {
    if cl.config.Target == "" {
        return nil, ErrInvalidTarget
    }
    client := cl.HttpClient()
    url := fmt.Sprintf("%s/users/%s/rooms", cl.config.Target, cl.config.Username)
    if resp, err := client.Get(url); err != nil {
        return nil, err
    } else if resp.StatusCode == http.StatusOK {
        var rooms []RoomData
        if data, err := ioutil.ReadAll(resp.Body); err != nil {
            return nil, err
        } else if err := json.Unmarshal(data, &rooms); err != nil {
            return nil, err
        } else {
            return rooms, nil
        }
    } else {
        return nil, fmt.Errorf("Error Status: %d", resp.StatusCode)
    }
}

func (cl *client) RoomRemove(name string) error {
    if cl.config.Target == "" {
        return ErrInvalidTarget
    }
    client := cl.HttpClient()
    url := fmt.Sprintf("%s/users/%s/rooms/%s", cl.config.Target, cl.config.Username, name)
    if req, err := http.NewRequest("DELETE", url, nil); err != nil {
        return err
    } else {
        if resp, err := client.Do(req); err != nil {
            return err
        } else if resp.StatusCode == http.StatusOK {
            return nil
        } else {
            return fmt.Errorf("Error Status: %d", resp.StatusCode)
        }
    }
    return nil
}

func (cl *client) RoomRename(currentName string, newName string) error {
    if cl.config.Target == "" {
        return ErrInvalidTarget
    }
    if js, err := json.Marshal(RoomBody{Name:newName}); err != nil {
        return err //barely impossible
    } else {
        client := cl.HttpClient()
        url := fmt.Sprintf("%s/users/%s/rooms/%s", cl.config.Target, cl.config.Username, currentName)
        if req, err := http.NewRequest("PATCH", url, bytes.NewReader(js)); err != nil {
            return err
        } else {
            req.Header.Set("Content-Type", JsonContentType)
            if resp, err := client.Do(req); err != nil {
                return err
            } else if resp.StatusCode == http.StatusOK {
                return nil
            } else {
                return fmt.Errorf("Error Status: %d", resp.StatusCode)
            }
        }
    }

    return nil
}

type SignUp struct {
    Username             string `json:"username"`
    Password             string `json:"password"`
    PasswordConfirmation string `json:"passwordConfirmation"`
}

func (cl *client) SignUp(username, password, passwordConfirmation string) error {
    if cl.config.Target == "" {
        return ErrInvalidTarget
    }
    body := SignUp{
        Username:username,
        Password:password,
        PasswordConfirmation: passwordConfirmation,
    }
    if js, err := json.Marshal(body); err != nil {
        return err //barely impossible
    } else {
        client := cl.HttpClient()
        url := fmt.Sprintf("%s/users", cl.config.Target)
        fmt.Println(url)
        if resp, err := client.Post(url, JsonContentType, bytes.NewReader(js)); err != nil {
            return err
        } else if resp.StatusCode == http.StatusOK {
            return nil
        } else {
            return fmt.Errorf("Error Status: %d", resp.StatusCode)
        }
    }
    return nil
}







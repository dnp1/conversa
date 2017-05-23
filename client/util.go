package client

import (
    "time"
    "net/http"
    "errors"
    "io"
    "io/ioutil"
    "encoding/json"
    "bytes"
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
    ResponseBody
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
type ResponseBody struct {
    Message string `json:"message"`
    Data json.RawMessage `json:"data"`
}


func ReadJSON(body io.Reader, refToData interface{}) error {
    if data, err := ioutil.ReadAll(body); err != nil {
        return err
    } else if err := json.Unmarshal(data, refToData); err != nil {
        return err
    }
    return nil
}

func ReadResponseBody(body io.ReadCloser) (*ResponseBody, Error) {
    defer body.Close()
    var respBody ResponseBody

    if err := ReadJSON(body, &respBody); err != nil {
        return nil, newServer(err)
    }
    return &respBody, nil
}

func JSONReader(data interface{}) (io.Reader, Error) {
    if js, err := json.Marshal(data); err != nil {
        return nil, newFatal(err) //barely impossible
    } else {
        return bytes.NewReader(js), nil
    }
}
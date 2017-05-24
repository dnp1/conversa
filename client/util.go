package client

import (
    "errors"
    "io"
    "io/ioutil"
    "encoding/json"
    "bytes"
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

type ResponseBody struct {
    Message string `json:"message"`
    Data    json.RawMessage `json:"data"`
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
    const kb = 1024
    limitedReader := io.LimitReader(body, 500 * kb)
    if err := ReadJSON(limitedReader, &respBody); err != nil {
        return nil, newServerError(err)
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

func IsOkCode(code int) bool {
    return code >= 200 && code < 300
}

func IsServerErrorCode(code int) bool {
    return code >= 500 && code < 600
}

func IsClientErrorCode(code int) bool {
    return code >= 400 && code < 500
}

func HttpError(body ResponseBody, code int) Error {
    if IsOkCode(code) {
        return nil
    } else if IsServerErrorCode(code) {
        err := errors.New(body.Message)
        return newServerError(err)
    } else if IsClientErrorCode(code) {
        err := errors.New(body.Message)
        return newServerError(err)
    } else {
        err := errors.New(body.Message)
        return newError(err)
    }

}

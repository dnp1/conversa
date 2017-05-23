package client

import (
    "time"
    "net/http"
    "io"
)

type requester struct {
    target string
}

func (r *requester) NewRequest(method, endpoint string, body io.Reader) (*http.Request, Error) {
    path := r.target + endpoint
    if req, err := http.NewRequest(method, path, body); err != nil {
        newFatal(err)
        return nil
    } else {
        if body != nil {
            req.Header.Set("Content-Type", "application/json")
        }
        return req
    }
}

func (r *requester) Do(req *http.Request, jar http.CookieJar) (*http.Response, Error) {
    c := &http.Client{
        Transport:transport,
        Timeout: time.Second * 15,
        Jar: jar,
    }
    resp, err := c.Do(req);
    return resp, newError(err)
}

func (r *requester) Request(method, endpoint string, body io.Reader, jar http.CookieJar) (*http.Response, Error) {
    if req, err := r.NewRequest(method, endpoint, body); err != nil {
        return nil, err
    } else {
        return r.Do(req, jar)
    }
}

func (r *requester) SimpleRequest(method, endpoint string, body io.Reader, jar http.CookieJar) (int, *ResponseBody, Error) {
    if req, err := r.NewRequest(method, endpoint, body); err != nil {
        return nil, err
    } else if resp, err := r.Do(req, jar); err != nil {
        return 0, nil, err
    } else if body, err := ReadResponseBody(resp.Body); err != nil {
        return 0, nil, err
    } else {
        return resp.StatusCode, body, nil
    }
}

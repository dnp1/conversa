package client

import (
    "time"
    "net/http"
    "io"
    "net/http/cookiejar"
)

type Requester interface {
    NewRequest(method, endpoint string, body io.Reader) (*http.Request)
    Do(req *http.Request, jar http.CookieJar) (*http.Response, error)
    Request(method, endpoint string, body io.Reader, jar http.CookieJar) (*http.Response, error)
    SimpleRequest(method, endpoint string, body io.Reader, jar http.CookieJar) (code int, response *ResponseBody, err Error)
}

type requester struct {
    target string
    client *http.Client
}

func newRequester(target string, cookies  []*http.Cookie) (Requester, Error) {
    var jar, err = cookiejar.New(nil)
    if err != nil {
        return nil, newFatal(err)
    }

    if cookies != nil {
        jar.SetCookies(target, cookies)
    }

    return &requester{
        target:target,
        client: &http.Client{
            Transport: &http.Transport{
                TLSHandshakeTimeout:10 * time.Second,
                ResponseHeaderTimeout: 10 * time.Second,
                MaxIdleConnsPerHost: 4,
                MaxIdleConns: 16,
                MaxResponseHeaderBytes: 4096,
                IdleConnTimeout: 2 * time.Minute,
                ExpectContinueTimeout: 30 * time.Second,
            },
            Timeout: 45 * time.Second,
            Jar: jar,
        },
    }, nil
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

func (r *requester) Do(req *http.Request) (*http.Response, Error) {
    resp, err := r.client.Do(req);
    return resp, newError(err)
}

func (r *requester) Request(method, endpoint string, body io.Reader) (*http.Response, Error) {
    if req, err := r.NewRequest(method, endpoint, body); err != nil {
        return nil, err
    } else {
        return r.Do(req)
    }
}

func (r *requester) SimpleRequest(method, endpoint string, body io.Reader) (int, *ResponseBody, Error) {
    if req, err := r.NewRequest(method, endpoint, body); err != nil {
        return nil, err
    } else if resp, err := r.Do(req); err != nil {
        return 0, nil, err
    } else if body, err := ReadResponseBody(resp.Body); err != nil {
        return 0, nil, err
    } else {
        return resp.StatusCode, body, nil
    }
}

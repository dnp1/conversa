package requester

import (
    "time"
    "net/http"
    "io"
    "net/http/cookiejar"
    "net/url"
    "github.com/dnp1/conversa/client/errors"
    "encoding/json"
    "bytes"
    "io/ioutil"
    "fmt"
)

type req struct {
    urlTarget url.URL
    client    *http.Client
}

func New(target string, cookies []*http.Cookie) (*req, errors.Error) {
    urlTarget, err := url.ParseRequestURI(target)
    if err != nil {
        return nil, errors.Unexpected(err)
    }

    jar, err := cookiejar.New(nil)
    if err != nil {
        return nil, errors.Unexpected(err)
    }

    if cookies != nil {
        jar.SetCookies(urlTarget, cookies)
    }

    return &req{
        urlTarget: *urlTarget,
        client: &http.Client{
            Transport: &http.Transport{
                TLSHandshakeTimeout:    15 * time.Second,
                ResponseHeaderTimeout:  30 * time.Second,
                MaxIdleConnsPerHost:    1000,
                MaxIdleConns:           4096,
                MaxResponseHeaderBytes: 1024,
                IdleConnTimeout:        2 * time.Minute,
                ExpectContinueTimeout:  30 * time.Second,
            },
            Timeout: 90 * time.Second,
            Jar:     jar,
        },
    }, nil
}

func (r *req) Request(method, path string, body, respDest interface{}) errors.Error {
    var (
        endpoint   = r.urlTarget
        bodyReader io.Reader
    )
    endpoint.Path = path
    if js, err := json.Marshal(body); err != nil {
        return errors.Unexpected(err)
    } else {
        bodyReader = bytes.NewReader(js)
    }
    req, err := http.NewRequest(method, endpoint.String(), bodyReader);
    if err != nil {
        return errors.Unexpected(err)
    }
    req.Header.Set("Content-Type", "application/json")
    if resp, err := r.client.Do(req); err != nil {
        return errors.Unexpected(err)
    } else {
        defer resp.Body.Close()
        if responseJs, err := ioutil.ReadAll(resp.Body); err != nil {
            return errors.Unexpected(err)
        } else if err := json.Unmarshal(responseJs, respDest); err != nil {
            return errors.Server(fmt.Errorf("Non json response (Code: %d)", resp.StatusCode))
        } else if err := errors.FromHttpStatus(resp.StatusCode, fmt.Sprintf("Error on request. %+v", respDest)); err != nil {
            return err
        }
    }
    return nil
}

func (r *req) Cookies() []*http.Cookie {
    var urlPath = r.urlTarget
    return r.client.Jar.Cookies(&urlPath)
}

package requester

import (
    "bytes"
    "bufio"
    "io"
    "net/http"
    "github.com/dnp1/conversa/client/errors"
    "fmt"
)

//SSE name constants
const (
    eventName = "event"
    dataName  = "data"
)

//Sse is a go representation of an http server-sent event
type Sse struct {
    Type string
    Data io.Reader
}

func (r *req) NotifySSE(path string, evCh chan<- *Sse, errCh chan <- errors.Error)  {
    var uri = r.urlTarget
    uri.Path = path
    req, err := http.NewRequest(http.MethodGet, uri.String(), nil)
    if err != nil {
        errCh <- errors.Unexpected(err)
        return
    }
    req.Header.Set("Accept", "text/event-stream")
    res, err := r.client.Do(req)
    if err != nil {
        err := fmt.Errorf("error performing request for %s: %v", uri, err)
        errCh <- errors.Unexpected(err)
        return
    }

    br := bufio.NewReader(res.Body)
    defer res.Body.Close()

    delimiter := []byte{':', ' '}

    var currEvent *Sse

    for {
        bs, err := br.ReadBytes('\n')

        if err != nil && err != io.EOF {
            errCh <- errors.Unexpected(err)
            return
        }

        if len(bs) < 2 {
            continue
        }

        spl := bytes.Split(bs, delimiter)

        if len(spl) < 2 {
            continue
        }

        currEvent = &Sse{}
        switch string(spl[0]) {
        case eventName:
            currEvent.Type = string(bytes.TrimSpace(spl[1]))
        case dataName:
            currEvent.Data = bytes.NewBuffer(bytes.TrimSpace(spl[1]))
            evCh <- currEvent
        }
        if err == io.EOF {
            break
        }
    }
}

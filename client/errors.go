package client

type baseError struct {
    fatal  bool
    client bool
    server bool
}

func (e *baseError) Server() bool {
    return e.client
}

func (e *baseError) Fatal() bool {
    return e.fatal
}

func (e *baseError) Client() bool {
    return e.client
}


type Error interface {
    error
    //Fatal means something are wrong in library it self
    Fatal() bool
    //Server means something went wrong at server
    Server() bool
    //Client means something on the request didn't matched the expected
    Client() bool
}

type serverError struct {
    error
    baseError
}

func newServerError(err error) Error {
    if err == nil {
        return nil
    }
    return &serverError{
        error:err,
        baseError{server:true},
    }
}

type fatalError struct {
    error
}

func newFatal(err error) Error {
    if err == nil {
        return nil
    }
    return &fatalError{
        error:err,
        baseError{fatal:true},
    }

}

type ordinaryError struct {
    error
    baseError
}

func newError(err error) Error {
    if err == nil {
        return nil
    }
    return &ordinaryError{
        error:err,
    }
}


type clientError struct {
    error
    baseError
}


func newClientError(err error) Error {
    if err == nil {
        return nil
    }
    return &clientError{
        error:err,
        baseError{client:true},
    }
}
package client

type Error interface {
    error
    Fatal() bool
    Server() bool
}

type serverError struct {
    error
}

func newServer(err error) Error {
    if err == nil {
        return nil
    }
    return &serverError{error:err}
}

func (e *serverError) Server() bool {
    return true
}

func (e *serverError) Fatal() bool {
    return false
}


type fatalError struct {
    error
}

func newFatal(err error) Error {
    if err == nil {
        return nil
    }
    return &fatalError{error:err}
}

func (e *fatalError) Server() bool {
    return false
}

func (e *fatalError) Fatal() bool {
    return true
}

type ordinaryError struct {
    error
}


func newError(err error) Error {
    if err == nil {
        return nil
    }
    return &ordinaryError{error:err}
}

func (e *ordinaryError) Server() bool {
    return false
}

func (e *ordinaryError) Fatal() bool {
    return false
}
package channel

import (
    "github.com/dnp1/conversa/server/data/message"
    "io"
    "github.com/dnp1/conversa/server/handlers"
)

func New(model Model) *handler {
    return &handler{
        model: model,
    }
}

type Model interface {
    Listen(username, roomName string) <-chan *message.EventData
    StopListening(username, roomName string, ch <-chan *message.EventData)
}

type handler struct {
    model Model
}

func (handler *handler) Listen(context handlers.ChannelContext) {
    var (
        username = context.Param("user")
        roomName = context.Param("room")
    )

    listener := handler.model.Listen(username, roomName)
    defer handler.model.StopListening(username, roomName, listener)
    context.Stream(
        func(w io.Writer) bool {
            event := <-listener
            context.SSEvent(event.Event, event.Data)
            return true

        })
}

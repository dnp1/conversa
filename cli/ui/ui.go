package ui

import (
    "fmt"
    "log"
    "github.com/jroimartin/gocui"
    "github.com/dnp1/conversa/client/errors"
    "time"
)

type Actions interface {
    MessageCreate(user, room, content string) errors.Error
}

type ChatUi struct {
    Username     string
    RoomName     string
    gui          *gocui.Gui
    viewMessages *gocui.View
    Actions Actions
}

func (ui *ChatUi) Init() {
    g, err := gocui.NewGui(gocui.Output256)
    g.Cursor = true
    g.Mouse = true
    ui.gui = g
    if err != nil {
        log.Panicln(err)
    }
    defer g.Close()

    g.SetManagerFunc(ui.layout)

    if err := ui.keyBindings(g); err != nil {
        log.Panicln(err)
    }

    if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
        log.Panicln(err)
    }
}

func (ui *ChatUi) ReceiveMessage(user string, message string, editedOn time.Time) {
    ui.gui.Execute(func(*gocui.Gui) error {
        fmt.Fprintf(ui.viewMessages, "(%s) %s:\n\t\t%s\n", editedOn.String(), user, message)
        return nil
    })
}

func (ui *ChatUi) layout(g *gocui.Gui) error {
    maxX, maxY := g.Size()
    if v, err := g.SetView("chat", 0, 0, maxX, maxY-3); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Editable = true
        v.Autoscroll = true
        v.Title = fmt.Sprintf("chat - %s/%s ", ui.Username, ui.RoomName)
        ui.viewMessages = v
    }
    if v, err := g.SetView("input", 0, maxY-3, maxX, maxY); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = "press <enter> to send message"
        v.Editable = true
        v.Wrap = true
        v.Highlight = true
        v.Editor = gocui.EditorFunc(ui.Editor)

        g.SetCurrentView("input")
    }
    return nil
}

func (ui *ChatUi) keyBindings(g *gocui.Gui) error {
    err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
        return gocui.ErrQuit
    })
    if err != nil {
        return err
    }

    if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone,
        func(g *gocui.Gui, v *gocui.View) error {
            scrollView(ui.viewMessages, -1)
            return nil
        }); err != nil {
        return err
    }

    if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone,
        func(g *gocui.Gui, v *gocui.View) error {
            scrollView(ui.viewMessages, 1)
            return nil
        }); err != nil {
        return err
    }

    if err := g.SetKeybinding("", gocui.KeyCtrlA, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
        v.Autoscroll = true
        return nil
    }); err != nil {
        return err
    }

    return nil
}

func scrollView(v *gocui.View, dy int) error {
    if v != nil {
        v.Autoscroll = false
        ox, oy := v.Origin()
        if err := v.SetOrigin(ox, oy+dy); err != nil {
            if dy > 0 {
                v.Autoscroll = true
            }
            return err
        }
    }
    return nil
}


func (ui *ChatUi) Editor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier)  {
    switch key {
    case gocui.KeyEnter:
        ui.Actions.MessageCreate(ui.Username, ui.RoomName, v.ViewBuffer())
        v.Clear()
    default:
        gocui.DefaultEditor.Edit(v, key, ch, mod)
    }
}


func (ui *ChatUi) Close(err errors.Error) {
    ui.gui.Close()
}



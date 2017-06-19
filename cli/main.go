package main

import (
    "os"
    "sort"
    "github.com/urfave/cli"
    "fmt"
    "syscall"
    "golang.org/x/crypto/ssh/terminal"
    "log"
    "bufio"
    "strings"
    "os/user"
    "path"
    "io/ioutil"
    "github.com/dnp1/conversa/client"
    "github.com/dnp1/conversa/client/errors"
    "github.com/dnp1/conversa/cli/ui"
)

func getConfigFilePath() string {
    if val, ok := os.LookupEnv("CONVERSA_CLIENT_SESSION"); ok {
        return val
    } else {
        usr, err := user.Current()
        if err != nil {
            log.Fatalln(err)
        }
        return path.Join(usr.HomeDir, ".conversa")
    }
}

func handleConf() ([]byte, error) {
    var configPath string = getConfigFilePath()
    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        return nil, nil
    } else if err != nil {
        return nil, err
    } else if data, err := ioutil.ReadFile(configPath); err != nil {
        return nil, err
    } else {
        return data, nil
    }
}

func readPassword(prompt string) string {
    cleanPrompt(prompt)
    bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
    if err != nil {
        log.Fatalf("Unexpected err: %s", err)
    }
    return string(bytePassword)
}

func cleanPrompt(prompt string) {
    prompt = strings.TrimSpace(prompt)
    fmt.Print(strings.TrimSuffix(prompt, ":") + ": ")
}

func readString(prompt string) string {
    cleanPrompt(prompt)
    r := bufio.NewReader(os.Stdin)
    if str, err := r.ReadString('\n'); err != nil {
        log.Fatalln(err)
        return ""
    } else {
        return strings.TrimSpace(str)
    }
}

type Api interface {
    Credentials() []byte
    Login(username, password string) errors.Error
    Logout() errors.Error
    RoomCreate(name string) errors.Error
    RoomList() ([]client.RoomItem, errors.Error)
    RoomRemove(name string) errors.Error
    SignUp(username, password, passwordConfirmation string) errors.Error
    MessageCreate(user, room, content string) errors.Error
    Listen(user, room string, ch chan<- *client.Message, errCh chan<- errors.Error)
}

func main() {
    var api Api
    if target, ok := os.LookupEnv("CONVERSA_TARGET_URL"); !ok {
        log.Println("You should provide CONVERSA_TARGET_URL env variable!!!")
        return
    } else if data, err := handleConf(); err != nil {
        log.Println(err)
        return
    } else {
        cl, err := client.New(target, data)
        if err != nil {
            log.Println(err)
            return
        }
        api = cl
    }

    app := cli.NewApp()

    app.Version = "1.0.0"
    app.Authors = []cli.Author{
        {
            Name: "Danilo A.N.P",
        },
    }
    app.EnableBashCompletion = true

    app.Commands = []cli.Command{
        {
            Name: "login",
            Flags: []cli.Flag{
                cli.StringFlag{
                    Name:  "user, u",
                    Usage: "-u `username`",
                },
                cli.StringFlag{
                    Name:  "password, p",
                    Usage: "-p  `password`",
                },
            },
            Usage: "Login in server",
            Action: func(c *cli.Context) error {
                var (
                    username string = c.String("user")
                    password string = c.String("password")
                )
                if username == "" {
                    username = readString("User")
                }
                if password == "" {
                    password = readPassword("Password")
                    fmt.Println()
                }
                if err := api.Login(username, password); err != nil {
                    return cli.NewMultiError(err)
                } else if ioutil.WriteFile(getConfigFilePath(), api.Credentials(), os.ModePerm); err != nil {
                    log.Println("Could not remove session file", err)
                } else {
                    fmt.Println()
                    fmt.Println("Successfully logged in!")
                }
                return nil
            },
        },
        {
            Name:  "logout",
            Usage: "Logout from server",
            Action: func(c *cli.Context) error {
                if err := api.Logout(); err != nil {
                    return cli.NewMultiError(err)
                } else {
                    if os.Remove(getConfigFilePath()); err != nil {
                        fmt.Println()
                        log.Println("Could not remove session file", err)
                    }
                    fmt.Println("Successfully logged out!")
                }
                return nil
            },
        },
        {
            Name:  "sign-up",
            Usage: "Create user account",
            Flags: []cli.Flag{
                cli.StringFlag{
                    Name:  "user, u",
                    Usage: "-u `username`",
                },
                cli.StringFlag{
                    Name:  "password, p",
                    Usage: "-p `password`",
                },
                cli.StringFlag{
                    Name:  "confirmation, c",
                    Usage: "-c `confirmation`",
                },
            },
            Action: func(c *cli.Context) error {
                var (
                    username             string = c.String("user")
                    password             string = c.String("password")
                    passwordConfirmation string = c.String("password")
                )

                if username == "" {
                    username = readString("User")
                }
                if password == "" {
                    password = readPassword("Password")
                    fmt.Println()
                }
                if passwordConfirmation == "" {
                    passwordConfirmation = readPassword("Confirm Password")
                    fmt.Println()
                }
                if err := api.SignUp(username, password, passwordConfirmation); err != nil {
                    fmt.Println()
                    return cli.NewMultiError(err)
                }
                return nil
            },
        },
        {
            Name:  "room-list",
            Usage: "[-u `username`]",
            Flags: []cli.Flag{
                cli.StringFlag{
                    Name:  "user, u",
                    Usage: "-u `username`",
                },
            },
            Action: func(c *cli.Context) error {
                var err errors.Error
                var data []client.RoomItem

                if data, err = api.RoomList(); err != nil {
                    if err.Authentication() {
                        fmt.Println("You must be signned in to perform this!")
                    } else {
                        fmt.Println(err)
                    }
                    return cli.NewMultiError(err)
                }
                fmt.Printf("[Rooms: %d]\n", len(data))
                username := c.String("user")
                for _, room := range data {
                    if username != "" && room.Username != username {
                        continue
                    }
                    fmt.Printf("\t%s .. %s\n", room.Username, room.Name)
                }
                return nil
            },
        },
        {
            Name:  "room-create",
            Usage: "-n name",
            Flags: []cli.Flag{
                cli.StringFlag{
                    Name:  "name, n",
                    Usage: "-n [name]",
                },
            },
            Action: func(c *cli.Context) error {
                var name = c.String("name")
                if name == "" {
                    name = readString("Room")
                }
                if err := api.RoomCreate(name); err != nil {
                    if err.Authentication() {
                        fmt.Println("You must be signned in to perform this!")
                    } else {
                        fmt.Println(err)
                    }

                }
                return nil
            },
        },
        {
            Name:  "room-remove",
            Usage: "-n `name`",
            Flags: []cli.Flag{
                cli.StringFlag{
                    Name:  "name, n",
                    Usage: "-n `name`",
                },
            },
            Action: func(c *cli.Context) error {
                var name = c.String("name")
                if name == "" {
                    name = readString("Room name")
                }

                if err := api.RoomRemove(name); err != nil {
                    if err.Authentication() {
                        fmt.Println("You must be signned in to perform this!")
                    } else {
                        fmt.Println(err)
                    }

                }

                return nil
            },
        },
        {
            Name:  "join-room",
            Usage: "-u username -n roomname",
            Flags: []cli.Flag{
                cli.StringFlag{
                    Name:  "name, n",
                    Usage: "-n `name`",
                },
                cli.StringFlag{
                    Name:  "user, u",
                    Usage: "-u `username`",
                },
            },
            Action: func(c *cli.Context) error {
                var username = c.String("user")
                if username == "" {
                    username = readString("Username")
                }
                var name = c.String("name")
                if name == "" {
                    name = readString("Room name")
                }

                i := ui.ChatUi{
                    Username: username,
                    RoomName: name,
                    Actions:  api,
                }

                ch := make(chan *client.Message)
                errCh := make(chan errors.Error)
                go api.Listen(username, name, ch, errCh)

                go func() {
                    for {
                        select {
                        case msg, ok := <-ch:
                             if !ok {
                                return
                            }

                            i.ReceiveMessage(msg.OwnerUsername, msg.Content, msg.EditionDatetime)


                        case err, ok := <-errCh:
                            if ok {
                                i.Close(err)
                                log.Println(err)
                            }
                            return
                        }
                    }
                }()

                i.Init()

                return nil
            },
        },
    }

    sort.Sort(cli.FlagsByName(app.Flags))
    sort.Sort(cli.CommandsByName(app.Commands))
    app.Run(os.Args)
}

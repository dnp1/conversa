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
    "github.com/dnp1/conversa/conversa-cli/cl"
)

func handleConf() *cl.Configuration {
    var configPath string
    if val, ok := os.LookupEnv("CONVERSA_CONFIG_FILE_PATH"); ok {
        configPath = val
    } else {
        usr, err := user.Current()
        if err != nil {
            log.Fatalln(err)
        }
        configPath = path.Join(usr.HomeDir, ".conversa")
    }
    if conf, err := cl.LoadFromPath(configPath); err != nil {
        conf = &cl.Configuration{Path:configPath}
        if err := conf.Create(); err != nil {
            log.Fatalln(err)
        } else {
            return conf
        }
    } else {
        return conf
    }
    return nil
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

func main() {
    conf := handleConf()
    client, err := conf.BuildClient()
    if err != nil {
        log.Fatalln(err)
    }
    app := cli.NewApp()

    app.Version = "1.0.0"
    app.Authors = []cli.Author{
        {
            Name:  "Danilo",
        },
    }
    app.EnableBashCompletion = true

    app.Commands = []cli.Command{
        {
            Name:    "target-set",
            Usage:   "--address `http[s]://(URL:PORT)`",

            Flags:[]cli.Flag{
                cli.StringFlag{
                    Name:  "address, a",
                    Usage: "Pass a url beginnign with http:// or https://",

                },
            },
            Action:  func(c *cli.Context) error {
                var address string
                if address = c.String("address"); address == "" {
                    address = readString("Please enter the url target")
                }
                conf.SetTarget(address)

                return nil
            },
        },
        {
            Name:    "login",
            Flags:[]cli.Flag{
                cli.StringFlag{
                    Name:  "user, u",
                    Usage: "-u `username`",
                },
                cli.StringFlag{
                    Name:  "password, p",
                    Usage: "-p  `password`",
                },
            },
            Usage:   "Login in server",
            Action:  func(c *cli.Context) error {
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
                if err := client.Login(username, password); err != nil {
                    return cli.NewMultiError(err)
                }
                return nil
            },
        },
        {
            Name:    "logout",
            Usage:   "Logout from server",
            Action:  func(c *cli.Context) error {
                if err := client.Logout(); err != nil {
                    return cli.NewMultiError(err)
                } else {
                    fmt.Println("Successfully logged out!")
                }
                return nil
            },
        },
        {
            Name:    "sign-up",
            Usage:   "Create user account",
            Flags:[]cli.Flag{
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
            Action:  func(c *cli.Context) error {
                var (
                    username string = c.String("user")
                    password string = c.String("password")
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
                if err := client.SignUp(username, password, passwordConfirmation); err != nil {
                    return cli.NewMultiError(err)
                }
                return nil
            },
        },
        {
            Name: "room-list",
            Usage: "[-u `username`]",
            Flags:[]cli.Flag{
                cli.StringFlag{
                    Name:  "user, u",
                    Usage: "-u `username`",
                },
            },
            Action:  func(c *cli.Context) error {
                var err error
                var data []cl.RoomData
                if username := c.String("user"); username != "" {
                    if data, err = client.RoomsUserList(username); err != nil {
                        return cli.NewMultiError(err)
                    }
                } else {
                    if data, err = client.RoomList(); err != nil {
                        return cli.NewMultiError(err)
                    }
                }

                fmt.Printf("[Rooms: %d]\n", len(data))
                for _, room := range data {
                    fmt.Printf("\t%s .. %s\n", room.Username, room.Name)
                }

                return nil
            },
        },
        {
            Name: "room-create",
            Usage: "-n name",
            Flags:[]cli.Flag{
                cli.StringFlag{
                    Name:  "name, n",
                    Usage: "-n [name]",
                },
            },
            Action:  func(c *cli.Context) error {
                var name = c.String("name")
                if name == "" {
                    name = readString("Room")
                }
                if err := client.RoomCreate(name); err != nil {
                    cli.NewMultiError(err)
                }
                return nil
            },
        },
        {
            Name: "room-remove",
            Usage: "-n `name`",
            Flags:[]cli.Flag{
                cli.StringFlag{
                    Name:  "name, n",
                    Usage: "-n `name`",
                },
            },
            Action:  func(c *cli.Context) error {
                var name = c.String("name")
                if name == "" {
                    name = readString("Room name")
                }

                if err := client.RoomRemove(name); err != nil {
                    cli.NewMultiError(err)
                }

                return nil
            },
        },
    }

    sort.Sort(cli.FlagsByName(app.Flags))
    sort.Sort(cli.CommandsByName(app.Commands))
    app.Run(os.Args)
}
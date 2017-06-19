package main

import (
    "log"
    "fmt"
    "time"
    "net/http"
    "os"
    _ "github.com/lib/pq"
    "database/sql"
    "github.com/dnp1/conversa/server/controller"
    authenticationHandlers "github.com/dnp1/conversa/server/handlers/authentication"
    roomHandlers "github.com/dnp1/conversa/server/handlers/room"
    userHandlers "github.com/dnp1/conversa/server/handlers/user"
    sessionHandlers "github.com/dnp1/conversa/server/handlers/session"
    sessionModel "github.com/dnp1/conversa/server/model/session"
    messageHandlers "github.com/dnp1/conversa/server/handlers/message"
    userModel "github.com/dnp1/conversa/server/model/user"
    "golang.org/x/crypto/bcrypt"
    roomModel "github.com/dnp1/conversa/server/model/room"
    messageModel "github.com/dnp1/conversa/server/model/message"
    "github.com/dnp1/conversa/server/handlers/channel"
)

func env(key string, defaultVal string) string {
    if val, exists := os.LookupEnv(key); exists {
        return val;
    } else {
        return defaultVal
    }
}

func main() {
    db, err := sql.Open("postgres", os.Getenv("CONVERSA_DB_CONN_STR"))
    if err != nil {
        log.Fatalln(err)
    }
    if err := db.Ping(); err != nil {
        log.Fatalln(err)
    }

    authenticationCookieName := env("CONVERSA_AUTHENTICATION_COOKIE_NAME", "Authentication")
    session := sessionModel.New(db)
    user := userModel.New(db, bcrypt.DefaultCost)
    room := roomModel.New(db)
    message := messageModel.New(db)

    router := controller.New(&controller.Handlers{
        Authentication: authenticationHandlers.New(authenticationCookieName, session),
        Message: messageHandlers.New(message),
        Session: sessionHandlers.New(session),
        Room: roomHandlers.New(room),
        User: userHandlers.New(user),
        Channel: channel.New(message),
    })

    host := env("HOST", "0.0.0.0")
    port := env("PORT", "8080")
    srv := &http.Server{
        Addr:              fmt.Sprintf("%s:%s", host, port),
        ReadTimeout:       60 * time.Second,
        WriteTimeout:      60 * time.Second,
        ReadHeaderTimeout: 10 * time.Second,
        MaxHeaderBytes:    1 << 11,
        Handler:           router,
    }

    if err := srv.ListenAndServe(); err != nil {
        log.Printf("Error when lister server %s", err)
    }
}

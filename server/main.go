package main

import (
    "log"
    "fmt"
    "time"
    "net/http"
    "os"
    "github.com/dnp1/conversa/server/controller"
    _ "github.com/lib/pq"
    "database/sql"
    "github.com/dnp1/conversa/server/model/room"
    "github.com/dnp1/conversa/server/model/session"
    "github.com/dnp1/conversa/server/model/user"
    "github.com/dnp1/conversa/server/model/message"
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

    host := env("HOST", "0.0.0.0")
    port := env("PORT", "5001")
    router := controller.RouterBuilder{
        Session: session.Builder{DB:db}.Build(),
        User: user.Builder{DB:db}.Build(),
        Room: room.Builder{DB:db}.Build(),
        Message: message.Builder{DB:db}.Build(),
    }.Build()

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

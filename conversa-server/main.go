package main

import (
    "log"
    "fmt"
    "time"
    "net/http"
    "os"
    "github.com/dnp1/conversa/conversa-server/server"
    _ "github.com/lib/pq"
    "database/sql"
)

func env(key string, defaultVal string) string {
    if val, exists := os.LookupEnv(key); exists {
        return val;
    } else {
        return defaultVal
    }
}

//inject dependencies here
func init() {

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
    srv := &http.Server{
        Addr: fmt.Sprintf("%s:%s", host, port),
        ReadTimeout: 60 * time.Second,
        WriteTimeout: 60 * time.Second,
        ReadHeaderTimeout: 10 * time.Second,
        MaxHeaderBytes: 1 << 11,
        Handler: server.NewRouter(db),
    }

    if err := srv.ListenAndServe(); err != nil {
        log.Printf("Error when lister server %s", err)
    }
}

package main


import "log"

//inject dependencies here
func init() {

}

func main() {
    srv := newServer()
    if err := srv.ListenAndServe(); err != nil {
        log.Printf("Error when lister server %s", err)
    }
}

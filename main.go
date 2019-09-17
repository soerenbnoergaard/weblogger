package main

import (
    "fmt"
    "log"
    "net/http"
)

const (
    TOKEN = "1155234"
)

func handler(w http.ResponseWriter, r *http.Request) {
    keys := r.URL.Query()

    if r.Method == "POST" {
        token, ok := keys["token"]
        if !ok {
            // log.Printf("No token given in request: %+v\n", r.URL)
            // return
        }
        if token[0] != TOKEN {
            log.Printf("Incorrect token: %s\n", token[0])
            return
        }

        data, ok := keys["data"]
        if !ok {
            // log.Printf("No token given in request: %+v\n", r.URL)
            // return
        }

        fmt.Fprintf(w, "Method = %s\n", r.Method)
        fmt.Fprintf(w, "Token = %s\n", token)
        fmt.Fprintf(w, "Data = %s\n", data)

    } else if r.Method == "GET" {
        fmt.Fprintf(w, "Method = %s\n", r.Method)
    }

}

func main() {
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

package main

import (
    "log"
    "net/http"
    "os"
)

func main() {
    dir := "public"
    if _, err := os.Stat(dir); os.IsNotExist(err) {
        dir = `e:\Projects\golang\massage-web\public`
    }
    fs := http.FileServer(http.Dir(dir))
    http.Handle("/", fs)
    log.Printf("Serving %s on http://localhost:8000\n", dir)
    log.Fatal(http.ListenAndServe(":8000", nil))
}

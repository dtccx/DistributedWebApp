package main

import (
  "net/http"
  "log"
)

func main() {
  fs := http.FileServer(http.Dir("static"))
  http.Handle("/", fs)

  log.Println("Listening...")
  http.ListenAndServe(":8080", nil)
}

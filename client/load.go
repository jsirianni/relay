package main

import (
    "os"
    "fmt"
    "time"
    "bytes"
    "net/http"
)


func main() {
    go fullSend()
    go fullSend()
    go fullSend()
    go fullSend()
    go fullSend()
    go fullSend()
    go fullSend()
    go fullSend()
    go fullSend()
    go fullSend()
    go fullSend()
    go fullSend()
    go fullSend()
    go fullSend()
    go fullSend()

    for {
        time.Sleep(time.Second * 3)
    }
}

func fullSend() {
    payload := []byte(`{"text":"xyz"}`)
    //c := 0
    for {
        url := "https://relay-test4.duckdns.org/message"
        req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        req.Header.Set("x-relay-api-key", "73733f48-2953-4ecc-a36c-c32782ca5ce2")
        req.Header.Set("Content-Type", "application/json")

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        defer resp.Body.Close()

    //    c = c + 1
    //    fmt.Println(c)
    }
}

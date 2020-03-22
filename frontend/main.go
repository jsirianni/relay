package main

import (
    "os"
    "fmt"
    "flag"
    "net/http"

    "github.com/jsirianni/relay/common"

    "cloud.google.com/go/pubsub"
    "github.com/gorilla/mux"
)

type IncomingRequest struct {
    Text string `json:"text"`
}

var p common.PubSubConf
var port string
var topic string

const apiKeyHeader   = "x-relay-api-key"
const invalidIPError = "ip address is not valid, failed to parse"
const missingAPIKeyHeader = "request did not include api key header"
const errTopicExists = "Resource already exists in the project"


func init() {
    flag.StringVar(&port, "port", "8080", "server http port")
    flag.StringVar(&topic, "topic", "", "pubsub topic to publish messages to")
    flag.Parse()
}

func main() {
    var err error
    p, err = common.Init()
    if err != nil {
        fmt.Fprint(os.Stderr, err.Error())
        os.Exit(1)
    }

    if topic == "" {
        p.Log.Error("topic must be set")
        os.Exit(1)
    }

    p.Topic = p.Client.Topic(topic)
    if err := server(); err != nil {
        p.Log.Error(err)
        p.Topic.Stop()
        os.Exit(1)
    }
    p.Topic.Stop()
    os.Exit(0)
}

func server() error {
    r := mux.NewRouter()
    r.HandleFunc("/message", handleMessage).Methods("POST")
    r.HandleFunc("/status", status).Methods("GET")
    p.Log.Trace("starting frontend relay server on port " + port)
    return http.ListenAndServe(":" + port, r)
}

func status(resp http.ResponseWriter, req *http.Request) {
    resp.WriteHeader(http.StatusOK)
}

func handleMessage(resp http.ResponseWriter, req *http.Request) {
    messageBytes, err := parseMessage(req)
    if err != nil {
        p.Log.Error(err)
        resp.WriteHeader(http.StatusInternalServerError)
        return
    }

    if err := send(messageBytes); err != nil {
        p.Log.Error(err)
        resp.WriteHeader(http.StatusInternalServerError)
        return
    }
    resp.WriteHeader(http.StatusOK)
}

func send(payload []byte) error {
    id, err := p.Topic.Publish(p.CTX, &pubsub.Message{Data: payload}).Get(p.CTX)
    if err != nil {
        return err
    }
    p.Log.Info("published message: " + id )
    return nil
}

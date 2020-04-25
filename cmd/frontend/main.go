package main

import (
    "os"
    "fmt"
    "flag"
    "net/http"

    "github.com/jsirianni/relay/internal/queue"
    "github.com/jsirianni/relay/internal/util/logger"
    "github.com/jsirianni/relay/internal/auth"
    "github.com/jsirianni/relay/internal/auth/gcpdatastore"
    "github.com/jsirianni/relay/internal/util/env"

    "github.com/gorilla/mux"
)

type Frontend struct {
    ProjectID string
    Queue     queue.Queue
    Auth      auth.Auth
    Log       logger.Logger
}

var front Frontend

var (
    port string
    topic string
)

const (
    apiKeyHeader   = "x-relay-api-key"
    invalidIPError = "ip address is not valid, failed to parse"
    missingAPIKeyHeader = "request did not include api key header"
    errTopicExists = "Resource already exists in the project"
)

func init() {
    flag.StringVar(&port, "port", "8080", "server http port")
    flag.StringVar(&topic, "topic", "", "pubsub topic to publish messages to")
    flag.Parse()

    if topic == "" {
        panic("topic must be set")
    }
}

func (f *Frontend) Init(topicName string) error {
    var err error

    f.ProjectID, err = env.ProjectID()
    if err != nil {
        return err
    }

    logLevel, err := env.LogLevel()
    if err != nil {
        return err
    }
    if logLevel == "" {
        logLevel = logger.InfoLVL
    }
    if err := front.Log.Configure(logLevel); err != nil {
        return err
    }

    f.Auth, err = gcpdatastore.New(f.ProjectID)
    if err != nil {
        return err
    }

    f.Queue, err = queue.New("google", topicName, f.Log)
    if err != nil {
        return err
    }

    return err
}

func main() {
    if err := front.Init(topic); err != nil {
        fmt.Fprint(os.Stderr, err.Error())
        os.Exit(1)
    }

    defer front.Queue.Stop()
    if err := server(); err != nil {
        front.Log.Error(err)
        os.Exit(1)
    }
    os.Exit(0)
}

func server() error {
    r := mux.NewRouter()
    r.HandleFunc("/message", handleMessage).Methods("POST")
    r.HandleFunc("/status", status).Methods("GET")
    front.Log.Info("starting frontend relay server on port " + port)
    front.Log.Info("using message topic: " + front.Queue.TopicName())
    return http.ListenAndServe(":" + port, r)
}

func status(resp http.ResponseWriter, req *http.Request) {
    if front.Log.Level() == logger.TraceLVL {
        addr, err := parseAddress(req)
        if err != nil {
            front.Log.Error(err)
        } else {
            front.Log.Trace("healthcheck from " + addr)
        }
    }
    resp.WriteHeader(http.StatusOK)
}

func handleMessage(resp http.ResponseWriter, req *http.Request) {
    addr, err := parseAddress(req)
    if err != nil {
        addr = "<could not parse ip>"
    }

    apiKey, err := parseAPIKey(req)
    if err != nil {
        front.Log.Error(err)
        resp.WriteHeader(http.StatusNetworkAuthenticationRequired)
        return
    }

    validAPIKey, err := front.Auth.ValidAPIKey(apiKey)
    if err != nil {
        front.Log.Error(err)
        resp.WriteHeader(http.StatusInternalServerError)
        return
    }
    if validAPIKey != true {
        front.Log.Trace("invalid api key " + apiKey + " from " + addr)
        resp.WriteHeader(http.StatusNetworkAuthenticationRequired)
        return
    }

    payload, err := parseMessage(req)
    if err != nil {
        front.Log.Error(err)
        resp.WriteHeader(http.StatusInternalServerError)
        return
    }

    if err := front.Queue.Publish(payload); err != nil {
        front.Log.Error(err)
        resp.WriteHeader(http.StatusInternalServerError)
        return
    }
    resp.WriteHeader(http.StatusOK)
}

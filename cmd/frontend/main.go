package main

import (
    "os"
    "fmt"
    "flag"
    "net/http"

    "github.com/jsirianni/relay/internal/queue/google"
    "github.com/jsirianni/relay/internal/util/logger"
    "github.com/jsirianni/relay/internal/auth"
    "github.com/jsirianni/relay/internal/auth/gcpdatastore"
    "github.com/jsirianni/relay/internal/env"

    "cloud.google.com/go/pubsub"
    "github.com/gorilla/mux"
)

type Frontend struct {
    ProjectID string
    PubSub    google.PubSubConf
    Auth      auth.Auth
    Log       logger.Logger
}

type IncomingRequest struct {
    Text string `json:"text"`
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
}

func (f *Frontend) Init() error {
    var err error

    f.ProjectID, err = env.ENVProjectID()
    if err != nil {
        return err
    }

    logLevel, err := env.ENVLogLevel()
    if err != nil {
        return err
    }

    if err := front.Log.Configure(logLevel); err != nil {
        return err
    }

    f.Auth, err = gcpdatastore.New(f.ProjectID)
    if err != nil {
        return err
    }

    f.PubSub, err = google.Init()
    if err != nil {
        return err
    }

    return err
}

func main() {
    if err := front.Init(); err != nil {
        fmt.Fprint(os.Stderr, err.Error())
        os.Exit(1)
    }

    if topic == "" {
        front.Log.Error("topic must be set")
        os.Exit(1)
    }

    front.PubSub.Topic = front.PubSub.Client.Topic(topic)
    if err := server(); err != nil {
        front.Log.Error(err)
        front.PubSub.Topic.Stop()
        os.Exit(1)
    }
    front.PubSub.Topic.Stop()
    os.Exit(0)
}

func server() error {
    r := mux.NewRouter()
    r.HandleFunc("/message", handleMessage).Methods("POST")
    r.HandleFunc("/status", status).Methods("GET")
    front.Log.Info("starting frontend relay server on port " + port)
    front.Log.Info("using message topic: " + front.PubSub.Topic.String())
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

    messageBytes, err := parseMessage(req)
    if err != nil {
        front.Log.Error(err)
        resp.WriteHeader(http.StatusInternalServerError)
        return
    }

    if err := send(messageBytes); err != nil {
        front.Log.Error(err)
        resp.WriteHeader(http.StatusInternalServerError)
        return
    }
    resp.WriteHeader(http.StatusOK)
}

func send(payload []byte) error {
    id, err := front.PubSub.Topic.Publish(front.PubSub.CTX, &pubsub.Message{Data: payload}).Get(front.PubSub.CTX)
    if err != nil {
        return err
    }
    front.Log.Info("published message: " + id )
    return nil
}

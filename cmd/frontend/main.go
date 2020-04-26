package main

import (
    "os"
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

var f Frontend

// globals set in init()
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
    logLevel, err := env.LogLevel()
    if err != nil {
        if !env.IsEnvNotSetError(err) {
            panic(err)
        }
        // set default value when environment is not set
        logLevel = logger.InfoLVL
    }
    if err := f.Log.Configure(logLevel); err != nil {
        panic(err)
    }

    port, err = env.FrontendPort()
    if err != nil {
        // panic if the error is something other than the
        // environment not being set. This indicates trouble with
        // the os package, which should not happen
        if !env.IsEnvNotSetError(err) {
            panic(err)
        }

        // set default value when environment is not set
        port = "8080"
    }

    topic, err = env.Topic()
    if err != nil {
        panic(err)
    }
}

func (f *Frontend) Init() error {
    var err error

    f.ProjectID, err = env.GoogleProjectID()
    if err != nil {
        return err
    }

    f.Auth, err = gcpdatastore.New(f.ProjectID)
    if err != nil {
        return err
    }

    f.Queue, err = queue.New("google", topic, f.Log)
    if err != nil {
        return err
    }

    return err
}

func main() {
    if err := f.Init(); err != nil {
        f.Log.Error(err)
        os.Exit(1)
    }

    defer f.Queue.Stop()
    if err := server(); err != nil {
        f.Log.Error(err)
        os.Exit(1)
    }
    os.Exit(0)
}

func server() error {
    r := mux.NewRouter()
    r.HandleFunc("/message", handleMessage).Methods("POST")
    r.HandleFunc("/status", status).Methods("GET")
    f.Log.Info("starting frontend relay server on port " + port)
    f.Log.Info("using message topic: " + f.Queue.TopicName())
    return http.ListenAndServe(":" + port, r)
}

func status(resp http.ResponseWriter, req *http.Request) {
    if f.Log.Level() == logger.TraceLVL {
        addr, err := parseAddress(req)
        if err != nil {
            f.Log.Error(err)
        } else {
            f.Log.Trace("healthcheck from " + addr)
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
        f.Log.Error(err)
        resp.WriteHeader(http.StatusNetworkAuthenticationRequired)
        return
    }

    validAPIKey, err := f.Auth.ValidAPIKey(apiKey)
    if err != nil {
        f.Log.Error(err)
        resp.WriteHeader(http.StatusInternalServerError)
        return
    }
    if validAPIKey != true {
        f.Log.Trace("invalid api key " + apiKey + " from " + addr)
        resp.WriteHeader(http.StatusNetworkAuthenticationRequired)
        return
    }

    payload, err := parseMessage(req)
    if err != nil {
        f.Log.Error(err)
        resp.WriteHeader(http.StatusInternalServerError)
        return
    }

    if err := f.Queue.Publish(payload); err != nil {
        f.Log.Error(err)
        resp.WriteHeader(http.StatusInternalServerError)
        return
    }
    resp.WriteHeader(http.StatusOK)
}

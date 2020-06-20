package main

import (
    "os"
    "net"
    "net/http"
    "strings"
    "context"
    "io/ioutil"

    "github.com/jsirianni/relay/internal/queue"
    "github.com/jsirianni/relay/internal/logger"
    "github.com/jsirianni/relay/internal/auth"
    "github.com/jsirianni/relay/internal/auth/gcpdatastore"
    "github.com/jsirianni/relay/internal/env"
    "github.com/jsirianni/relay/internal/message"

    "github.com/gorilla/mux"
    "github.com/google/uuid"
    "github.com/pkg/errors"
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
    r.HandleFunc("/message", messageHandler).Methods("POST")
    r.HandleFunc("/status", statusHandler).Methods("GET")
    f.Log.Info("starting frontend relay server on port " + port)
    f.Log.Info("using message topic: " + f.Queue.TopicName())
    return http.ListenAndServe(":" + port, r)
}

func statusHandler(resp http.ResponseWriter, req *http.Request) {
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

func messageHandler(resp http.ResponseWriter, req *http.Request) {
    m := message.New()
    m.SetTime()
    m.CTX = context.Background()

    addr, err := parseAddress(req)
    if err != nil {
        addr = "<could not parse ip>"
    }
    m.SetAddress(addr)

    apiKey, err := uuid.Parse(req.Header.Get(apiKeyHeader))
    if err != nil {
        f.Log.Trace(err)
        resp.WriteHeader(http.StatusNetworkAuthenticationRequired)
        return
    }
    m.SetAPIKey(apiKey)

    valid, err := f.Auth.ValidAPIKey(apiKey)
    if err != nil {
        f.Log.Error(err)
        resp.WriteHeader(http.StatusInternalServerError)
        return
    }
    if valid != true {
        f.Log.Trace("invalid api key " + apiKey.String() + " from " + addr)
        resp.WriteHeader(http.StatusNetworkAuthenticationRequired)
        return
    }



    p, err := ioutil.ReadAll(req.Body)
    if err != nil {
        f.Log.Error(err)
        resp.WriteHeader(http.StatusInternalServerError)
        return
    }

    if err := m.ParsePayload(p); err != nil {
        f.Log.Error(err)
        resp.WriteHeader(http.StatusInternalServerError)
        return
    }

    // safely log the message without the APIKey
    safeJson, err := m.BytesSafe()
    if err != nil {
        f.Log.Error(err)
    }
    f.Log.Info("new message: " + string(safeJson))

    // return the message as json
    payload, err := m.Bytes()
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

func parseAddress(req *http.Request) (string, error) {
    raw := strings.Split(req.RemoteAddr, ":")[0]
    if raw == "[" {
        raw = "127.0.0.1"
    }
    addr := net.ParseIP(raw)
    if addr == nil {
        f.Log.Trace(errors.Wrap(errors.New(invalidIPError), "failed to parse address from '" + raw + "'"))
        return "", errors.New(invalidIPError)
    }
    return addr.String(), nil
}

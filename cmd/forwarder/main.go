package main

import (
    "os"
    "sync"
    "flag"
    "encoding/json"

    "github.com/jsirianni/relay/internal/message"
    "github.com/jsirianni/relay/internal/queue"
    "github.com/jsirianni/relay/internal/queue/qmessage"
    "github.com/jsirianni/relay/internal/alert"
    "github.com/jsirianni/relay/internal/util/env"
    "github.com/jsirianni/relay/internal/util/logger"
)

type Forwarder struct {
    Alert       alert.Alert
    Queue       queue.Queue
    QueueType   string
    Log         logger.Logger
}

var f Forwarder
var queueType string
var subscription string

func init() {
    logLevel, err := env.LogLevel()
    if err != nil {
        panic(err)
    }
    if logLevel == "" {
        logLevel = logger.InfoLVL
    }
    if err := f.Log.Configure(logLevel); err != nil {
        panic(err)
    }

    flag.StringVar(&subscription, "subscription", "", "pubsub subscription to listen on")
    flag.StringVar(&queueType, "queue-type", "", "message queue type (defaults to Google Pubsub)")
    flag.Parse()

    if subscription == "" {
        panic("subscription must be set")
    }

    if queueType == "" {
        queueType = "google"
    }
}

func main() {
    var err error
    f.Queue, err = queue.New(queueType, subscription, f.Log)
    if err != nil {
        f.Log.Error(err)
        os.Exit(1)
    }

    f.Alert, err = initDest()
    if err != nil {
        f.Log.Error(err)
        os.Exit(1)
    }
    confBytes, err := f.Alert.Config()
    if err != nil {
        f.Log.Trace(err)
    } else {
        f.Log.Trace("destination configured with config: " + string(confBytes))
    }

    wg := sync.Mutex{}
    c := make(chan qmessage.Message)

    go f.Queue.Listen(subscription, c, &wg)
    for {
        m := <-c
        wg.Lock()
        if err := process(m.Payload); err != nil {
            f.Log.Trace(err.Error())
            m.ERR = err
        } else {
            m.ACK = true
        }
        c <- m
        wg.Unlock()
    }
}

func process(mRaw []byte) error {
    m := message.New()
    err := json.Unmarshal(mRaw, &m)
    if err != nil {
        return err
    }

    mSafe, err := m.BytesSafe()
    if err != nil {
        return err
    }
    f.Log.Info("new message: " + string(mSafe))

    if err := f.Alert.Message(m.Payload.Text); err != nil {
        return err
    }
    f.Log.Trace("message sent to destination '" + f.Alert.Type() + "'")
    return nil
}

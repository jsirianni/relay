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
    "github.com/jsirianni/relay/internal/env"
    "github.com/jsirianni/relay/internal/util/logger"
)

type Forwarder struct {
    Destination string
    Alert       alert.Alert
    Queue       queue.Queue
    Log         logger.Logger
}

var f Forwarder
//var destination alert.Alert
var subscription string

func init() {
    flag.StringVar(&subscription, "subscription", "", "pubsub subscription to listen on")
    flag.Parse()

    if subscription == "" {
        panic("subscription must be set")
    }

    logLevel, err := env.ENVLogLevel()
    if err != nil {
        panic(err)
    }
    if err := f.Log.Configure(logLevel); err != nil {
        panic(err)
    }
}

func main() {
    var err error
    f.Queue, err = queue.New("google", subscription, f.Log)
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

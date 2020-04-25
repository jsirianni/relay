package main

import (
    "os"
    "sync"
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

// globals set in init()
var (
    queueType string
    subscription string
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

    subscription, err = env.Subscription()
    if err != nil {
        panic(err)
    }

    queueType, err = env.QueueType()
    if err != nil {
        panic(err)
    }
}

func (f *Forwarder) Init() error {
    var err error
    
    f.Queue, err = queue.New(queueType, subscription, f.Log)
    if err != nil {
        return err
    }

    f.Alert, err = initDest()
    if err != nil {
        return err
    }
    confBytes, err := f.Alert.Config()
    if err != nil {
        return err
    } else {
        f.Log.Trace("destination configured with config: " + string(confBytes))
    }

    return nil
}

func main() {
    if err := f.Init(); err != nil {
        f.Log.Error(err)
        os.Exit(1)
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

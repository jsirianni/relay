package main

import (
    "fmt"
    "os"
    "flag"
    "context"
    "encoding/json"

    "github.com/jsirianni/relay/internal/queue/google"
    "github.com/jsirianni/relay/internal/message"
    "github.com/jsirianni/relay/internal/alert"

    "cloud.google.com/go/pubsub"
)

var p google.PubSubConf
var destination alert.Alert
var subscription string

func init() {
    flag.StringVar(&subscription, "subscription", "", "pubsub subscription to listen on")
    flag.Parse()
}

func main() {
    var err error
    p, err = google.Init()
    if err != nil {
        fmt.Fprint(os.Stderr, err.Error())
        os.Exit(1)
    }

    destination, err = initDest()
    if err != nil {
        p.Log.Error(err)
        os.Exit(1)
    }
    p.Log.Trace("destination configured")
    subscribe()
}

func subscribe() {
    p.Log.Info("starting listener on subscription: " + subscription)
    sub := p.Client.Subscription(subscription)
    cctx, _ := context.WithCancel(p.CTX)
    err := sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
        if err := process(msg.Data); err != nil {
            p.Log.Error(err)
        } else {
            msg.Ack()
            p.Log.Trace("ack sent for message " + msg.ID)
        }
    })
    if err != nil {
        p.Log.Error(err)
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
    p.Log.Info("new message: " + string(mSafe))

    if err := destination.Message(m.Text); err != nil {
        return err
    }
    p.Log.Trace("message sent without error")
    return nil
}
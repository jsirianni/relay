package google

import (
    "sync"
    "context"

    "github.com/jsirianni/relay/internal/env"
    "github.com/jsirianni/relay/internal/util/logger"
    "github.com/jsirianni/relay/internal/queue/qmessage"

    "cloud.google.com/go/pubsub"
    "github.com/pkg/errors"
)

type Google struct {
    ProjectID string

    CTX    context.Context
    Client *pubsub.Client
    Topic  *pubsub.Topic

    // TODO pass a logger to this instead of reading the environment
    Log    logger.Logger
}

func Init(topic string, l logger.Logger) (Google, error) {
    var (
        p        Google
        err      error
    )

    p.Log = l
    if p.Log.Configured() == false {
        return p, errors.New("logger is not configured")
    }

    p.ProjectID, err = env.ENVProjectID()
    if err != nil {
        return p, err
    }

    ctx := context.Background()
    p.CTX = ctx

    p.Client, err = pubsub.NewClient(p.CTX, p.ProjectID)
    if err != nil {
        return p, err
    }

    p.Topic = p.Client.Topic(topic)

    return p, nil
}

func (p Google) Publish(payload []byte) error {
    id, err := p.Topic.Publish(p.CTX, &pubsub.Message{Data: payload}).Get(p.CTX)
    if err != nil {
        return err
    }
    p.Log.Info("published message: " + id )
    return nil
}

func (p Google) Stop() {
    p.Topic.Stop()
}

func (p Google) TopicName() string {
    return p.Topic.String()
}

func (p Google) Listen(s string, c chan qmessage.Message, wg *sync.Mutex,) {
    // lock right away, before handling any messages
    wg.Lock()
    defer wg.Unlock()

    p.Log.Info("starting listener on subscription: " + s)
    sub := p.Client.Subscription(s)
    cctx, _ := context.WithCancel(p.CTX)
    err := sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
        // push the pubsub message through the channel and then
        // wait for it to be processed
        m := qmessage.Message{Payload: msg.Data}
        c <- m
        wg.Unlock()

        // wait for message to be processed
        for {
            m := <-c
            wg.Lock()
            if m.ACK == true {
                // TODO: does this return an error that
                // needs to be handled
                msg.Ack()
                return
            }
            if m.ERR != nil {
                p.Log.Error(m.ERR)
                return
            }
            // TODO: can this be handled better? What should be done here? We could re-publish or just
            // take the next message from pubsub?
            p.Log.Error(errors.New("qmessage was returned without and error and without ACK being set to true"))
            return
        }
    })
    // TODO: check if this is a complete failure, if so, panic. Make sure to document this
    // panic as it is not good practice to panic from a lib. We could return an error
    // but the intention here is to run this as a go routine
    if err != nil {
        p.Log.Error(err)
    }
}

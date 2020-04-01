package google

import (
    "context"

    "github.com/jsirianni/relay/internal/env"
    "github.com/jsirianni/relay/internal/util/logger"

    "cloud.google.com/go/pubsub"
)





type PubSubConf struct {
    ProjectID string

    CTX    context.Context
    Client *pubsub.Client
    Topic  *pubsub.Topic

    Log    logger.Logger
}

func Init() (PubSubConf, error) {
    var (
        p        PubSubConf
        logLevel string
        err      error
    )

    logLevel, err = env.ENVLogLevel()
    if err != nil {
        return p, err
    }
    if err := p.Log.Configure(logLevel); err != nil {
        return p, err
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
    return p, nil
}

package queue

import (
    "sync"

    "github.com/jsirianni/relay/internal/util/logger"
    "github.com/jsirianni/relay/internal/queue/qmessage"
    "github.com/jsirianni/relay/internal/queue/google"

    "github.com/pkg/errors"
)

const notSupportedERR = "specified message queue is not implemented"
const topicNotSetERR  = "topic is not set"

type Queue interface{
    Publish([]byte) error
    Listen(string, chan qmessage.Message, *sync.Mutex)

    Stop()
    TopicName() string
}

func New(queueType, topic string, l logger.Logger) (Queue, error) {
    if err := validate(topic, l); err != nil {
        return nil, err
    }

    if queueType == "google" {
        return google.Init(topic, l)
    }

    return nil, errors.New(notSupportedERR)
}

func validate(topic string, l logger.Logger) error {
    if topic == "" {
        return errors.New(topicNotSetERR)
    }

    if l.Configured() == false {
        return errors.New("logger is not configured")
    }

    return nil
}

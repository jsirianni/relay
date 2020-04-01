package message

import (
    "time"
    "encoding/json"

    "github.com/google/uuid"
    "github.com/pkg/errors"
)

type Message struct {
    APIKey uuid.UUID
    Text   string

    // UTC unix timestamp in nano seconds
    TimeStamp int64
    Address   string
}

func New() Message {
    return Message{}
}

func (m Message) Bytes() ([]byte, error) {
    return json.Marshal(m)
}

func (m Message) BytesSafe() ([]byte, error) {
    var err error
    newM := m
    emptyUUID := []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
    newM.APIKey, err = uuid.FromBytes(emptyUUID)
    if err != nil {
        return nil, err
    }
    return newM.Bytes()
}

func (m *Message) SetAPIKey(a string) (err error) {
    m.APIKey, err = uuid.Parse(a)
    return err
}

func (m *Message) SetText(t string) error {
    if t == "" {
        return errors.New("message text cannot be empty")
    }
    m.Text = t
    return nil
}

func (m *Message) SetTime() {
    m.TimeStamp = time.Now().UTC().UnixNano()
}

func (m *Message) SetAddress(addr string) {
    m.Address = addr
}

package message

import (
    "time"
    "encoding/json"

    "github.com/google/uuid"
)

type Message struct {
    APIKey uuid.UUID

    Payload struct {
        Text string `json:"text"`
        Type string `json:"type"`
    }

    // UTC unix timestamp in nano seconds
    TimeStamp int64
    Address   string
}

func New() Message {
    return Message{}
}

// Bytes returns the message as a json object
func (m Message) Bytes() ([]byte, error) {
    return json.Marshal(m)
}

// BytesSafe returns the mssage as a json object but with an
// empty APIKey, safe for logging
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

func (m *Message) SetTime() {
    m.TimeStamp = time.Now().UTC().UnixNano()
}

func (m *Message) SetAddress(addr string) {
    m.Address = addr
}

func (m *Message) ParsePayload(p []byte) error {
    return json.Unmarshal(p, &m.Payload)
}

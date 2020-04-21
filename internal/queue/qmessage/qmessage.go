package qmessage

// Message represents a payload returned by a type implementing
// the queue interface. ACK is set to true when the message
// is handled, and ERR is set to an error if the message
// was not able to be handled
type Message struct {
    Payload []byte
    ACK bool
    ERR error
}

package framework

// type based on byte slice
type Payload []byte

func (p Payload) String() string {
        return string(p)
}

func NewPayloadString(str string) Payload {
        return Payload([]byte(str))
}

func NewPayload(b []byte) Payload {
        return Payload(b)
}

type Message struct {
        senderPID int
        clock   Clock
        payload Payload
}

func (msg *Message) GetClock() Clock {
        return msg.clock
}

func (msg *Message)GetPayload() Payload {
        return msg.payload
}

func (msg *Message)SenderPID() int {
        return msg.senderPID
}

func NewMessage(senderPID int, clock Clock, payload Payload) *Message {
        return &Message{senderPID, clock, payload}
}

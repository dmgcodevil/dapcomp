package framework

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
        clock   ScalarTime
        payload Payload
}

func (msg Message) GetClock() ScalarTime {
        return msg.clock
}

func (msg Message)GetPayload() Payload {
        return msg.payload
}

func NewMessage(clock ScalarTime, payload Payload) Message {
        return Message{clock, payload}
}

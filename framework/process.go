package framework

import (
        "strconv"
        "fmt"
)

type Process struct {
        Pid     int
        st      *ScalarTime
        history []Message
        counter int
}

func (p *Process) LocalTime() ScalarTime {
        return *p.st
}

func NewProcess(pid int) *Process {
        return &Process{pid, NewScalarTime(), make([]Message, 0), 0}
}

func (p *Process) SendInternalMessage() {
        p.st.Update()
        p.counter++
        payload := "internal message" + strconv.Itoa(p.counter)
        msg := NewMessage(*p.st, NewPayloadString(payload))
        p.history = append(p.history, msg)
}

func (p *Process) SendMsg(payload Payload, c chan Message) {
        p.st.Update()
        msg := NewMessage(*p.st, payload)
        p.history = append(p.history, msg)
        c <- msg
}

func (p *Process)Receive(msg Message) {
        if ok, updated := p.st.Resolve(msg.GetClock()); ok == true {
                p.st = updated
                p.history = append(p.history, NewMessage(*updated, msg.GetPayload()))
        }
}

func (p *Process) PrintHistory() {
        fmt.Printf("pid: %d\n", p.Pid)
        for _, m := range p.history {
                msg := Message(m)
                fmt.Printf("message[payload: %s, scalar time: %d]\n", msg.GetPayload().String(), msg.GetClock().Value())
        }
}
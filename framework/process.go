package framework

import (
        "strconv"
        "fmt"
        "sync"
        "log"
        "os"
)

var mu sync.RWMutex

const MESSAGE_RECEIVED_FORMAT = "current local clock = [%s]; message clock = [%s]";

var (
        infoLog    *log.Logger = log.New(os.Stdout,
                "INFO: ",
                log.Ldate | log.Ltime | log.Lshortfile)
        warningLog *log.Logger = log.New(os.Stdout,
                "WARNING: ",
                log.Ldate | log.Ltime | log.Lshortfile)
        errorLog   *log.Logger = log.New(os.Stderr,
                "ERROR: ",
                log.Ldate | log.Ltime | log.Lshortfile)
)

type Process struct {
        Pid         int
        localTime   Clock
        history     []*Message
        counter     int
        listener    *Listener
        broadcaster *Broadcaster
        started     bool
}

func (p *Process) LocalTime() Clock {
        mu.RLock()
        defer mu.RUnlock()
        return p.localTime
}

func NewProcess(pid int, b *Broadcaster) (p *Process) {
        if l, e := b.Listen(pid); e == nil {
                p = &Process{pid, NewScalarTime(), make([]*Message, 0), 0, l, b, false}
        } else {
                panic(e)
        }

        return
}

func (p *Process) Stop() {
        mu.Lock()
        defer mu.Unlock()
        if (p.started) {
                p.listener.Close()
                p.started = false
        } else {
                panic(fmt.Sprintf("attempt to stop inactive process [%d]", p.Pid))
        }

}

func (p *Process) Start(wg *sync.WaitGroup) {
        mu.Lock()
        if !p.started {
                p.started = true;
                mu.Unlock()
                for {
                        m, more := <-p.listener.Ch
                        if (more) {
                                p.receive(m)
                        } else {
                                wg.Done()
                                return
                        }
                }
        } else {
                panic(fmt.Sprintf("Process [%d] is already running", p.Pid))
                mu.Unlock()
        }
}

func (p *Process) SendInternalMessage() {
        p.localTime = p.localTime.Update(p.Pid)
        p.counter++
        payload := "internal message: " + strconv.Itoa(p.counter)
        msg := NewMessage(p.Pid, p.localTime, NewPayloadString(payload))
        p.history = append(p.history, msg)
}

func (p *Process) SendMsg(payload Payload) {
        // R1: update process localTime before send
        p.localTime = p.localTime.Update(p.Pid)
        msg := NewMessage(p.Pid, p.localTime, payload)
        p.history = append(p.history, msg) // store message in process local storage, only for debug
        p.broadcaster.Send(msg)
}

func (p *Process)receive(msg *Message) {
        mu.Lock()
        defer mu.Unlock()
        inClock := msg.GetClock()
        if updated, ok := p.localTime.UpdateFrom(p.Pid, inClock); ok == true {
                infoLog.Print(
                        fmt.Sprintf("%s -> Received message. %s", p.String(), clockInfo(p, msg)))
                p.localTime = updated // update process local clock
                // create new message with updated clock value and deliver it to current process
                msgToDeliver := NewMessage(msg.senderPID, updated, msg.GetPayload())
                p.deliver(msgToDeliver)

        } else {
                warningLog.Print(fmt.Sprintf("%s -> Received stale message. %s", p.String(), clockInfo(p, msg)))
        }
}

func clockInfo(p *Process, msg *Message) string {
        return fmt.Sprintf(MESSAGE_RECEIVED_FORMAT, p.localTime.String(), msg.GetClock().String())
}

func (p *Process) deliver(msg *Message) {
        p.history = append(p.history, msg)
}

func (p *Process) String() string {
        return fmt.Sprintf("Process [%d]", p.Pid)
}

func (p *Process) PrintHistory() {
        infoLog.Println(fmt.Sprintf("%s: local clock = %s", p.String(), p.localTime.String()))
        infoLog.Println(fmt.Sprintf("%s: Messages", p.String()))
        for _, m := range p.history {
                c := m.GetClock()
                infoLog.Println(fmt.Sprintf("%s -> Message[payload: %s, clock: %s", p.String(), m.GetPayload().String(), c.String()))
        }
}
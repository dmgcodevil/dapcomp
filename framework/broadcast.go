package framework

import (
        "sync"
        "bytes"
        "fmt"
        "log"
        "errors"
)

var initListenersOnce sync.Once // lazy init



type listeners map[int]chan <- *Message

type Broadcaster struct {
        m         sync.Mutex
        listeners listeners
        capacity  int
        closed    bool
}

func (b *Broadcaster) initListeners() {
        b.listeners = make(listeners)
}

func New(n int) *Broadcaster {
        return &Broadcaster{capacity: n}
}

type Listener struct {
        Ch <-chan *Message
        b  *Broadcaster
        id int
}

func (b *Broadcaster) Send(m *Message) {
        b.m.Lock()
        defer b.m.Unlock()
        if b.closed {
                panic("broadcast: send after close")
        }
        for pid, ch := range b.listeners {
                if (pid != m.senderPID) {
                        ch <- m
                }
        }
}

func (b *Broadcaster) Close() {
        b.m.Lock()
        defer b.m.Unlock()
        b.closed = true
        for _, l := range b.listeners {
                close(l)
        }
}

func (b *Broadcaster) Listen(pid int) (l *Listener, e error) {
        b.m.Lock()
        defer b.m.Unlock()
        var buf bytes.Buffer
        logger := log.New(&buf, "Broadcaster: ", log.Lshortfile)
        initListenersOnce.Do(b.initListeners)
        l, e = nil, nil
        if b.closed {
                e = errors.New("broadcaster is closed")
                return
        }

        if (b.listeners[pid] != nil) {
                e = errors.New("process is alredy registerd")
                return
        }

        ch := make(chan *Message, b.capacity)
        logger.Printf("Prosses pid=%d has subscribed", pid)
        fmt.Print(&buf)
        b.listeners[pid] = ch

        l = &Listener{ch, b, pid}
        return
}

func (l *Listener) Close() {
        l.b.m.Lock()
        defer l.b.m.Unlock()
        ch := l.b.listeners[l.id]
        close(ch)
        delete(l.b.listeners, l.id)
}
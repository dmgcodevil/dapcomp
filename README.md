# dapcomp

## Logical Time

**Introduction**

**Framework**


### Scalar time

**Definition**
The scalar time representation was proposed by Lamport in 1978 as an attempt to totally order events in a distributed system.
In this approach both local process P(_i_) and global clocks squashed into one non-negative interger number C(_i_).

There are two rules used to update clock value:
* Rule 1: before executing send or internal event, process P(_i_) executes the following:
C(_i_) = C(_i_) + d       (d > 0)
Value of _d_ can be application dependent and have different value, however, typically _d_ is kept at 1 because it's enough to identify 
the time of each event uniquely. 
* Rule 2: Each message piggybacks the clock value of it's sender at sending time. When a process p(_i_) receives a message with timestamp 
C(_msg_), it executes the following actions:
    * C(_i_) = _max_(C(_i_), C(_msg_))
    * execute rule 1
    * deliver the message

Lets take a look at example:

```go
package main

import (
        "github.com/dmgcodevil/distributedprogramming/framework"
        "sync"
)

var broadcaster = &framework.Broadcaster{}

func createProcessAndStart(wg *sync.WaitGroup, pid int) (p *framework.Process) {
        wg.Add(1)
        p = framework.NewProcess(pid, broadcaster)
        go p.Start(wg)
        return
}

func startSimulation() {
        var wg sync.WaitGroup

        p1 := createProcessAndStart(&wg, 1)
        p2 := createProcessAndStart(&wg, 2)
        p3 := createProcessAndStart(&wg, 3)

        p1.SendMsg(framework.NewPayloadString("m1"))
        p2.SendMsg(framework.NewPayloadString("m2"))
        p3.SendMsg(framework.NewPayloadString("m3"))

        p1.Stop()
        p2.Stop()
        p3.Stop()

        wg.Wait()

        println("Current status of all processes")
        p1.PrintHistory()
        p2.PrintHistory()
        p3.PrintHistory()

}

func main() {
        startSimulation()
}
```

Output:

```
Broadcaster: broadcast.go:78: Prosses pid=1 has subscribed
Broadcaster: broadcast.go:78: Prosses pid=2 has subscribed
Broadcaster: broadcast.go:78: Prosses pid=3 has subscribed
INFO: 2016/08/06 21:44:17 process.go:107: Process [2] -> Received message. current local clock = [0]; message clock = [1]
INFO: 2016/08/06 21:44:17 process.go:107: Process [3] -> Received message. current local clock = [0]; message clock = [1]
INFO: 2016/08/06 21:44:17 process.go:107: Process [3] -> Received message. current local clock = [2]; message clock = [3]
INFO: 2016/08/06 21:44:17 process.go:107: Process [1] -> Received message. current local clock = [1]; message clock = [3]
INFO: 2016/08/06 21:44:17 process.go:107: Process [1] -> Received message. current local clock = [4]; message clock = [5]
INFO: 2016/08/06 21:44:17 process.go:107: Process [2] -> Received message. current local clock = [3]; message clock = [5]
INFO: 2016/08/06 21:44:17 process.go:131: Process [1]: local clock = 6
INFO: 2016/08/06 21:44:17 process.go:132: Process [1]: Messages
INFO: 2016/08/06 21:44:17 process.go:135: Process [1] -> Message[payload: m1, clock: 1
INFO: 2016/08/06 21:44:17 process.go:135: Process [1] -> Message[payload: m2, clock: 4
INFO: 2016/08/06 21:44:17 process.go:135: Process [1] -> Message[payload: m3, clock: 6
INFO: 2016/08/06 21:44:17 process.go:131: Process [2]: local clock = 6
INFO: 2016/08/06 21:44:17 process.go:132: Process [2]: Messages
INFO: 2016/08/06 21:44:17 process.go:135: Process [2] -> Message[payload: m1, clock: 2
INFO: 2016/08/06 21:44:17 process.go:135: Process [2] -> Message[payload: m2, clock: 3
INFO: 2016/08/06 21:44:17 process.go:135: Process [2] -> Message[payload: m3, clock: 6
INFO: 2016/08/06 21:44:17 process.go:131: Process [3]: local clock = 5
INFO: 2016/08/06 21:44:17 process.go:132: Process [3]: Messages
INFO: 2016/08/06 21:44:17 process.go:135: Process [3] -> Message[payload: m1, clock: 2
INFO: 2016/08/06 21:44:17 process.go:135: Process [3] -> Message[payload: m2, clock: 4
INFO: 2016/08/06 21:44:17 process.go:135: Process [3] -> Message[payload: m3, clock: 5
```

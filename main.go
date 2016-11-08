package main

import (
        "github.com/dmgcodevil/distributedprogramming/framework"
        "sync"
        "fmt"
)

var broadcaster = &framework.Broadcaster{}

type ClockType int

const (
        SCALAR ClockType = 1 + iota
        VECTOR
        SK_VECTOR // Singhal–Kshemkalyani’s differential technique
)

func createClock(clockType ClockType, numOfProcesses int) func() framework.Clock {
        switch clockType {
        case SCALAR:
                return func() framework.Clock { return framework.NewScalarTime() }
        case VECTOR:
                return func() framework.Clock { return framework.NewVectorClock(numOfProcesses) }
        default:
                panic(fmt.Sprintf("unsupported clock type: %d", clockType))
        }
}

func createProcessAndStart(wg *sync.WaitGroup, pid int, clockType ClockType, numOfProcesses int) (p *framework.Process) {
        wg.Add(1)
        p = framework.NewProcess(pid, broadcaster, createClock(clockType, numOfProcesses))
        go p.Start(wg)
        return
}

func startSimulation() {
        var wg sync.WaitGroup
        numOfProcesses := 3
        clockType := VECTOR
        p1 := createProcessAndStart(&wg, 0, clockType, numOfProcesses)
        p2 := createProcessAndStart(&wg, 1, clockType, numOfProcesses)
        p3 := createProcessAndStart(&wg, 2, clockType, numOfProcesses)

        p1.SendMsg(framework.NewPayloadString("m1"))
        p2.SendMsg(framework.NewPayloadString("m2"))
        p3.SendMsg(framework.NewPayloadString("m3"))
        p1.SendMsg(framework.NewPayloadString("m1"))

        p1.Stop()
        p2.Stop()
        p3.Stop()

        wg.Wait()

        println("messages per process")
        p1.PrintHistory()
        p2.PrintHistory()
        p3.PrintHistory()

        println("current clock")
        p1.PrintClock()
        p2.PrintClock()
        p3.PrintClock()

}

func main() {
        startSimulation()
}
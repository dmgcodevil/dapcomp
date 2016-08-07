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
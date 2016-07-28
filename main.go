package main

import (
        "github.com/dmgcodevil/distributedprogramming/framework"
)

func main() {
        // test create process
        ch := make(chan framework.Message)
        p := framework.NewProcess(1)
        p.SendInternalMessage()
        go p.SendMsg(framework.NewPayloadString("send message to external process"), ch)
        _ = <-ch
        p.PrintHistory()

        // create new msg
        staleMsg := framework.NewMessage(*framework.NewScalarTimeClock(2), framework.NewPayloadString("stale message"))
        newMsg := framework.NewMessage(*framework.NewScalarTimeClock(3), framework.NewPayloadString("new message"))
        p.Receive(staleMsg)
        p.Receive(newMsg)

        println("-----------------------------")
        p.PrintHistory()

}
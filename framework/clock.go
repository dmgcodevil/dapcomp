package framework

import "fmt"

const d = 1

type Clock interface {
        Update(pid int) Clock
        UpdateFrom(pid int, c Clock) (Clock, bool)
        String() string
}

func ToScalarTime(c Clock) ScalarTime {
        st, ok := c.(ScalarTime)
        if ok {
                return st
        }
        panic("failed to cast given clock to 'ScalarTime' type")
}

// Scalar Time
type ScalarTime struct {
        clock int
}

func (st ScalarTime) Value() int {
        return st.clock
}

func NewScalarTime() ScalarTime {
        return ScalarTime{0}
}

func NewScalarTimeClock(initial int) ScalarTime {
        return ScalarTime{initial}
}

func (st ScalarTime) Update(pid int) Clock {
        return ScalarTime{st.clock + d}
}

func (st ScalarTime) String() string {
        return fmt.Sprintf("%d", st.clock);
}

func (st ScalarTime) UpdateFrom(pid int, c Clock) (Clock, bool) {
        clock, ok := c.(ScalarTime)
        if ok {
                if (st.Value() < clock.Value()) {
                        return clock.Update(pid), true
                }
                return st, false
        } else {
                panic("given clock isn't instance of ScalarTime struct")
        }
}

// Vector Clock
type VectorClock struct {
        size  int
        clock []int
}

func (vc VectorClock) Update(pid int) Clock {
        return vc;
}

func (vc VectorClock) UpdateFrom(pid int, c Clock) (Clock, bool) {
        clock, ok := c.(VectorClock)
        if ok {
                println(clock.size) // todo implement
                return clock, true
        } else {
                panic("given clock isn't instance of VectorClock struct")
        }
}

func (vc VectorClock)Value() []int {
        return vc.clock
}

func (vc VectorClock) String() string {
        return fmt.Sprintf("%v", vc.clock);
}

func NewVectorClock(n int) VectorClock {
        return VectorClock{n, make([]int, n)}
}
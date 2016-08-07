package framework

import (
        "fmt"
)

const d = 1

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

func (st ScalarTime) UpdateAndGet() ScalarTime {
        return ScalarTime{st.clock + d}
}

func (st ScalarTime) String() string {
        return fmt.Sprintf("%d", st.clock);
}

func (st ScalarTime) Resolve(clock ScalarTime) (bool, ScalarTime) {
        if (st.Value() < clock.Value()) {
                return true, clock.UpdateAndGet()
        }
        return false, st
}
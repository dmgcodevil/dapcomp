package framework

const d = 1

type ScalarTime struct {
        clock int
}

func (st ScalarTime) Value() int {
        return st.clock
}

func NewScalarTime() *ScalarTime {
        return &ScalarTime{0}
}

func NewScalarTimeClock(clock int ) *ScalarTime {
        return &ScalarTime{clock}
}

func (st *ScalarTime) Update() *ScalarTime {
        st.clock = st.clock + d
        return st
}

func (st *ScalarTime) Resolve(delivered ScalarTime) (bool, *ScalarTime) {
        if (st.Value() < delivered.Value()) {
                st.clock = delivered.clock
                st.Update()
                return true, st
        }
        return false, st
}
package framework

import "testing"

type UpdateFromResult struct{ Clock; bool }

func TestScalarTime_Update(t *testing.T) {
        pid := 1
        given := NewScalarTime()
        res := given.Update(pid)

        st, ok := res.(ScalarTime)
        if (!ok) {
                t.Error("Update function should return value of 'ScalarTime' type")
        }
        if (1 != st.Value()) {
                t.Error("expected clock time is '1'")
        }
}

func TestScalarTime_UpdateFrom(t *testing.T) {

        pid := 1

        var fibTests = []struct {
                currentClock  Clock
                receivedClock Clock
                expected      UpdateFromResult
        }{
                {NewScalarTimeClock(1), NewScalarTimeClock(2), UpdateFromResult{NewScalarTimeClock(3), true}},
                {NewScalarTimeClock(2), NewScalarTimeClock(1), UpdateFromResult{NewScalarTimeClock(2), false}},

        }

        for _, tt := range fibTests {
                c, ok := tt.currentClock.UpdateFrom(pid, tt.receivedClock)
                if expectedResult := tt.expected.bool; ok != tt.expected.bool {
                        t.Errorf("expected (clock, %t) result of 'UpdateFrom' operation", expectedResult)
                };
                if expectedClock, updated := ToScalarTime(tt.expected.Clock), ToScalarTime(c); updated.Value() != expectedClock.Value() {
                        t.Errorf("expected clock value '%d', actual '%d'", expectedClock.Value(), updated.Value())
                }
        }

}
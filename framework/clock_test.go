package framework

import (
        "testing"
        "reflect"
)

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

func TestVectorClock_Update(t *testing.T) {
        pid := 0
        vc := NewVectorClock(3);
        updated := ToVectorClock(vc.Update(pid))

        expectedLocalTime := 1
        if (updated.clock[pid] != expectedLocalTime) {
                t.Errorf("expected process local time = %d", expectedLocalTime)
        }
}

func TestVectorClock_UpdateFrom(t *testing.T) {
        pid := 0
        i := 3
        localClock := NewVectorClock(i)
        msgClock := NewVectorClock(i);
        msgClock = ToVectorClock(msgClock.Update(1))

        c, ok := localClock.UpdateFrom(pid, msgClock)

        if (!ok) {
                t.Errorf("clock should be upddated");
                return
        }
        assertArrayEquals(t, []int{0, 1, 0}, ToVectorClock(c).clock)

}

func assertArrayEquals(t *testing.T, expected interface{}, actual interface{}) {
        validateArrayType := func(arr interface{}) {
                arrValue := reflect.ValueOf(arr)
                arrType := arrValue.Type()
                if arrType.Kind() != reflect.Array && arrType.Kind() != reflect.Slice {
                        panic("Array parameter's type is neither array nor slice.")
                }
        }
        validateArrayType(expected)
        validateArrayType(actual)

        expectedArrType, expectedArrayVal, expectedArrElemType := arrayTypeInfo(expected)
        actualArrType, actualArrayVal, actualArrElemType := arrayTypeInfo(actual)

        if (expectedArrType != actualArrType) {
                t.Errorf("different types of arrays, expected array type = %T, actual = %T", expectedArrType, actualArrType);
        }
        if (expectedArrElemType != actualArrElemType) {
                t.Errorf("different element types of arrays, expected elem type = %T, actual = %T", expectedArrElemType, actualArrElemType);
        }

        if (expectedArrayVal.Len() != actualArrayVal.Len()) {
                t.Errorf("different sizes of arrays, expected size = %d, actual size = %d", expectedArrayVal.Len(), actualArrayVal.Len());
        } else {
                for i := 0; i < expectedArrayVal.Len(); i++ {
                        actualVal := actualArrayVal.Index(i).Interface()
                        expectedVal := expectedArrayVal.Index(i).Interface()
                        if (expectedVal != actualVal) {
                                t.Errorf("different values at index = %d, expected = %v, actual = %v", i, expectedVal, actualVal);
                                return
                        }
                }

        }
}

func arrayTypeInfo(arr interface{}) (reflect.Type, reflect.Value, reflect.Type) {
        arrValue := reflect.ValueOf(arr)
        arrType := arrValue.Type()
        arrElemType := arrType.Elem()
        return arrType, arrValue, arrElemType
}

func assertPanic(t *testing.T, f func()) {
        defer func() {
                if r := recover(); r == nil {
                        t.Errorf("The code did not panic")
                }
        }()
        f()
}
package timeincreaser

import (
	"fmt"
	// "math/rand"

	"time"
)

type testStruct struct {
	data [2]int64
}

type testIncreaser struct {
}

func (t *testIncreaser) Max(param any) int64 {
	return 100
}

func (t *testIncreaser) Min(param any) int64 {
	return 0
}

func (t *testIncreaser) IncInterval(param any) int64 {
	return 1000
}

func (t *testIncreaser) IncCount(param any) float64 {
	return 0.5
}

func (t *testIncreaser) GetData(param any) (int64, int64) {
	p := param.(*testStruct)
	ts := p.data[0]
	count := p.data[1]
	if ts == 0 {
		ts = time.Now().Unix()
	}
	return count, ts
}

func (t *testIncreaser) SetData(count int64, ts int64, param any) {
	p := param.(*testStruct)
	p.data[0] = ts
	p.data[1] = count
}

func Example() {
	s := &testStruct{}
	inc := &testIncreaser{}
	now := time.Now().Unix()

	count, ts, addCount := IncreaserGet(inc, now, s)
	fmt.Println("init data:", count, addCount, ts == now)

	count, ts, _ = IncreaserSettle(inc, now, s)
	fmt.Println("after settle:", count, ts == now)
	fmt.Println("now data:", s.data[0] == now, s.data[1])

	now += 10000
	count, _, _ = IncreaserGet(inc, now, s)
	fmt.Println("after 10 seconds:", count)

	now += 500 * 1000
	count, _, _ = IncreaserGet(inc, now, s)
	fmt.Println("after 500 seconds, reach the max:", count)

	count, _, _ = IncreaserAdd(inc, 10, now, false, s)
	fmt.Println("add still works even after the max:", count, s.data[0] == now, s.data[1])

	_, ok, _ := IncreaserAdd(inc, -120, now, true, s)
	fmt.Println("remove will be failed:", ok)

	now += 100 * 1000
	count, _, _ = IncreaserAdd(inc, -100, now, false, s)
	fmt.Println("if the value reaches max, timestamp will reset after value changing from max(or more) to nonmax", count, s.data[0] == now, s.data[1])

	// Output:
	// init data: 0 0 true
	// after settle: 0 true
	// now data: true 0
	// after 10 seconds: 5
	// after 500 seconds, reach the max: 100
	// add still works even after the max: 110 true 110
	// remove will be failed: false
	// if the value reaches max, timestamp will reset after value changing from max(or more) to nonmax 10 true 10
}

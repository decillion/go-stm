package stm

import (
	"fmt"
	"sync"
)

func Example() {
	x := New(0)
	y := New(0)
	wg := sync.WaitGroup{}

	// Atomically increment x and y.
	inc := func(rec *TRec) interface{} {
		currX := rec.Load(x).(int)
		currY := rec.Load(y).(int)
		rec.Store(x, currX+1)
		rec.Store(y, currY+1)
		return nil
	}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			Atomically(inc)
			wg.Done()
		}()
	}
	wg.Wait()

	// Read values of x and y atomically.
	load := func(rec *TRec) interface{} {
		var curr [2]int
		curr[0] = rec.Load(x).(int)
		curr[1] = rec.Load(y).(int)
		return curr
	}
	curr := Atomically(load).([2]int)

	fmt.Printf("x = %v, y = %v", curr[0], curr[1])
	// Output:
	// x = 100, y = 100
}

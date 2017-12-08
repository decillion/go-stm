package stm

import (
	"fmt"
	"sync"
)

func Example() {
	x := New(0)
	y := New(0)

	// Define a transaction that increments x and y.
	inc := func(rec *TRec) {
		currX := rec.Load(x).(int)
		currY := rec.Load(y).(int)
		rec.Store(x, currX+1)
		rec.Store(y, currY+1)
	}

	// Run the transaction concurrently.
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			Atomically(inc)
			wg.Done()
		}()
	}
	wg.Wait()

	// Read values of x and y atomically.
	var currX, currY int // local variable
	load := func(rec *TRec) {
		currX = rec.Load(x).(int)
		currY = rec.Load(y).(int)
	}
	Atomically(load)

	fmt.Printf("x = %v, y = %v", currX, currY)
	// Output:
	// x = 100, y = 100
}

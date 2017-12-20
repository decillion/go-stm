package stm

import (
	"fmt"
	"sync"
	"testing"
)

func Example() {
	// There are two bank accounts of Alice and Bob.
	accountA := New(100)
	accountB := New(0)

	// Transfer 20 from Alice's account to Bob's one.
	transfer := func(rec *TRec) interface{} {
		currA := rec.Load(accountA).(int)
		currB := rec.Load(accountB).(int)
		rec.Store(accountA, currA-20)
		rec.Store(accountB, currB+20)
		return nil
	}
	Atomically(transfer)

	// Check the balance of accounts of Alice and Bob.
	inquiries := func(rec *TRec) interface{} {
		balance := make(map[*TVar]int)
		balance[accountA] = rec.Load(accountA).(int)
		balance[accountB] = rec.Load(accountB).(int)
		return balance
	}
	balance := Atomically(inquiries).(map[*TVar]int)

	fmt.Printf("The account of Alice holds %v.\nThe account of Bob holds %v.",
		balance[accountA], balance[accountB])
	// Output:
	// The account of Alice holds 80.
	// The account of Bob holds 20.
}

func TestIncrement(t *testing.T) {
	iter := 10000

	x := New(0)
	y := New(0)

	inc := func(rec *TRec) interface{} {
		rec.Store(x, rec.Load(x).(int)+1)
		rec.Store(y, rec.Load(y).(int)+1)
		return nil
	}

	read := func(rec *TRec) interface{} {
		var curr [2]int
		curr[0] = rec.Load(x).(int)
		curr[1] = rec.Load(y).(int)
		return curr
	}

	wg := sync.WaitGroup{}

	for i := 0; i < 2*iter; i++ {
		wg.Add(1)
		j := i
		go func() {
			if j%2 == 0 {
				Atomically(inc)
			} else {
				Atomically(read)
			}
			wg.Done()
		}()
	}

	wg.Wait()

	curr := Atomically(read).([2]int)
	currX, currY := curr[0], curr[1]

	if currX != iter || currY != iter {
		t.Errorf("want: (x, y) = (%v, %v) got: (x, y) = (%v, %v)",
			iter, iter, currX, currY)
	}
}

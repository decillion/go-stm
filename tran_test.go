package stm

import "fmt"

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

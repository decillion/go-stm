# go-stm

    import "github.com/decillion/go-stm"

Package stm is a software transactional memory implementation for Go, which is
based on the Transactional Locking II (TL2) proposed by Dice et al.
https://doi.org/10.1007/11864219_14

## Documents

Please see the [godoc page](https://godoc.org/github.com/decillion/go-stm) for further information.

## Example 

```go
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
```

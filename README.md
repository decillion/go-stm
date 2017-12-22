# go-stm

    import "github.com/decillion/go-stm"

Package stm is a software transactional memory implementation for Go, which is
based on the Transactional Locking II (TL2) proposed by Dice et al.
https://doi.org/10.1007/11864219_14

## Documents

See the [godoc page](https://godoc.org/github.com/decillion/go-stm) for further information.
Here is a [blog post](https://qiita.com/decillion/items/d5da905e28b54dc769cd) about go-stm (in Japanese).

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

## Benchmark

There exists an another [STM package](https://github.com/lukechampine/stm) written by lukechampine.
The lukechampine's package provides a richer interface but is less efficient than the present package.
Here is the result of a very simple benchmark, in which two transactional variables are atomically incremented or read.
The benchmark is taken at a DigitalOcean's High CPU Droplet with 32 cores.

    Benchmark_Read90Write10_decillion-2      	10000000	       230 ns/op
    Benchmark_Read90Write10_decillion-4      	10000000	       156 ns/op
    Benchmark_Read90Write10_decillion-8      	10000000	       144 ns/op
    Benchmark_Read90Write10_decillion-16       	10000000	       214 ns/op
    Benchmark_Read90Write10_decillion-32       	 5000000	       289 ns/op

    Benchmark_Read90Write10_lukechampine-2   	 2000000	       715 ns/op
    Benchmark_Read90Write10_lukechampine-4   	 2000000	       761 ns/op
    Benchmark_Read90Write10_lukechampine-8   	 2000000	       822 ns/op
    Benchmark_Read90Write10_lukechampine-16    	 2000000	       912 ns/op
    Benchmark_Read90Write10_lukechampine-32    	 2000000	       966 ns/op

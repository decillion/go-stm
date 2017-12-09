# go-stm

    import "github.com/decillion/go-stm"

Package stm is a software transactional memory implementation for Go, which is
based on the Transactional Locking II (TL2) proposed by Dice et al.
https://doi.org/10.1007/11864219_14

## Example 

```go
x := stm.New(0)
y := stm.New(0)
wg := sync.WaitGroup{}

// Atomically increment x and y.
inc := func(rec *stm.TRec) {
    currX := rec.Load(x).(int)
    currY := rec.Load(y).(int)
    rec.Store(x, currX+1)
    rec.Store(y, currY+1)
}
for i := 0; i < 100; i++ {
    wg.Add(1)
    go func() {
        stm.Atomically(inc)
        wg.Done()
    }()
}
wg.Wait()


// Read values of x and y atomically.
var currX, currY int // local variables
load := func(rec *stm.TRec) {
	currX = rec.Load(x).(int)
	currY = rec.Load(y).(int)
}
stm.Atomically(load)

fmt.Printf("x = %v, y = %v", currX, currY)
// Output:
// x = 100, y = 100
```

## Documents

Please see the [godoc page](https://godoc.org/github.com/decillion/go-stm) for further information.
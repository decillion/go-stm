// Package stm is a software transactional memory implementation for Go, which
// is based on the Transactional Locking II (TL2) proposed by Dice et al.
// https://doi.org/10.1007/11864219_14
package stm

import "sync/atomic"

type clock uint64

func (c *clock) incAndFetch() (t uint64) {
	return atomic.AddUint64((*uint64)(c), 1)
}

func (c *clock) sampleClock() (t uint64) {
	return atomic.LoadUint64((*uint64)(c))
}

var globalClock clock

// TVar is a transactional variable.
type TVar struct {
	lock  versionedLock
	value atomic.Value
}

// New creates a transactional variable and initialize it by the given value v.
func New(v interface{}) (x *TVar) {
	x = &TVar{}
	x.value.Store(v)
	return x
}

// TRec is an auxiliary data type for recording meta data of a transaction.
type TRec struct {
	aborted      bool
	readVersion  uint64
	writeVersion uint64
	readSet      map[*TVar]struct{}
	writeSet     map[*TVar]interface{}
}

// Load returns the value of the transactional variable x.
func (rec *TRec) Load(x *TVar) (v interface{}) {
	if v, ok := rec.writeSet[x]; ok {
		return v // No validation
	}
	if rec.readSet == nil {
		rec.readSet = make(map[*TVar]struct{})
	}
	rec.readSet[x] = struct{}{}
	_, pre := x.lock.sampleLock()
	v = x.value.Load()
	locked, post := x.lock.sampleLock()
	if locked || pre != post || post > rec.readVersion {
		rec.aborted = true
	}
	return v
}

// Store sets the value of the transactional variable x to the given value v.
func (rec *TRec) Store(x *TVar, v interface{}) {
	if rec.writeSet == nil {
		rec.writeSet = make(map[*TVar]interface{})
	}
	rec.writeSet[x] = v
}

// Atomically executes the given transaction tx atomically. The transaction tx
// should not contain non-transactional shared-variables.
func Atomically(tx func(rec *TRec) interface{}) interface{} {
RETRY:

	rec := &TRec{
		readVersion: globalClock.sampleClock(),
	}

	v := tx(rec) // speculative execution
	if rec.aborted {
		goto RETRY
	}

	if len(rec.writeSet) == 0 {
		return v // A fast path for read-only transaction
	}

	lockedSet := make(map[*TVar]struct{})
	for x := range rec.writeSet {
		if !x.lock.tryLock() {
			for x := range lockedSet {
				x.lock.unlock()
			}
			goto RETRY
		}
		lockedSet[x] = struct{}{}
	}

	rec.writeVersion = globalClock.incAndFetch()

	// Validate the elements of the read-set.
	if rec.writeVersion != rec.readVersion+1 {
		for x := range rec.readSet {
			locked, version := x.lock.sampleLock()
			_, lockedByMe := lockedSet[x]
			if (!lockedByMe && locked) || version > rec.readVersion {
				for y := range lockedSet {
					y.lock.unlock()
				}
				goto RETRY
			}
		}
	}

	for x, val := range rec.writeSet {
		x.value.Store(val)
		x.lock.unlockAndUpdate(rec.writeVersion)
	}
	return v
}

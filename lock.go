package stm

import "sync/atomic"

// versionedLock virtually maintains two fields, a lock bit and a version
// number. The first bit of uint64 is used to hold the lock bit and the rest
// is used to hold the version number. The zero value of versionedLock is
// (locked, version) = (false, 0).
type versionedLock uint64

func encode(locked bool, version uint64) (encoded uint64) {
	if (version >> 63) == 1 {
		panic("stm: version number exceeds 2^63-1.")
	}
	if locked {
		return (1 << 63) | version
	}
	return version
}

func decode(encoded uint64) (locked bool, version uint64) {
	if (encoded >> 63) == 1 {
		locked = true
	}
	version = (1<<63 - 1) & encoded
	return
}

func (lock *versionedLock) sampleLock() (locked bool, version uint64) {
	return decode(atomic.LoadUint64((*uint64)(lock)))
}

func (lock *versionedLock) tryLock() (ok bool) {
	old := atomic.LoadUint64((*uint64)(lock))
	locked, version := decode(old)
	if locked {
		return false
	}
	new := encode(true, version)
	return atomic.CompareAndSwapUint64((*uint64)(lock), old, new)
}

func (lock *versionedLock) unlock() {
	locked, version := lock.sampleLock()
	if !locked {
		panic("stm: unlock of unlocked versioned-lock.")
	}
	new := encode(false, version)
	atomic.StoreUint64((*uint64)(lock), new)
}

func (lock *versionedLock) unlockAndUpdate(version uint64) {
	locked, _ := lock.sampleLock()
	if !locked {
		panic("stm: unlock of unlocked versioned-lock.")
	}
	new := encode(false, version)
	atomic.StoreUint64((*uint64)(lock), new)
}

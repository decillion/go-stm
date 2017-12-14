package stm

import "testing"

const (
	minVersion = 0
	midVersion = 1 << 32
)

func helperEncode(locked bool, version uint64) (lock *versionedLock) {
	new := (versionedLock)(encode(locked, version))
	return &new
}

func helperDecode(lock *versionedLock) (locked bool, version uint64) {
	return decode((uint64)(*lock))
}

func TestTryLock(t *testing.T) {
	cases := []struct {
		locked  bool
		version uint64
	}{
		{false, minVersion},
		{false, midVersion},
		{true, minVersion},
		{true, midVersion},
	}

	for _, tc := range cases {
		lock := helperEncode(tc.locked, tc.version)

		ok := lock.tryLock()
		if ok == tc.locked {
			t.Errorf("got: ok = %v, want: ok = %v", ok, !ok)
		}

		locked, version := helperDecode(lock)
		if locked != true || version != tc.version {
			t.Errorf("got: (%v, %v) want: (%v, %v)", locked, version, true, tc.version)
		}
	}
}

func TestUnlock(t *testing.T) {
	cases := []struct {
		locked  bool
		version uint64
	}{
		{true, minVersion},
		{true, midVersion},
	}

	for _, tc := range cases {
		lock := helperEncode(tc.locked, tc.version)
		lock.unlock()

		locked, version := helperDecode(lock)
		if locked != false || version != tc.version {
			t.Errorf("got: (%v, %v) want: (%v, %v)", locked, version, false, tc.version)
		}
	}
}

func TestUnlockAndUpdate(t *testing.T) {
	cases := []struct {
		locked  bool
		version uint64
		input   uint64
	}{
		{true, minVersion, minVersion + 1},
		{true, midVersion, midVersion + 1},
	}

	for _, tc := range cases {
		lock := helperEncode(tc.locked, tc.version)
		lock.unlockAndUpdate(tc.input)

		locked, version := helperDecode(lock)
		if locked != false || version != tc.input {
			t.Errorf("got: (%v, %v) want: (%v, %v)", locked, version, false, tc.input)
		}
	}
}

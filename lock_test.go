package stm

import "testing"

const (
	minVersion = 0
	midVersion = 1 << 30
	maxVersion = 1<<63 - 1
)

func TestEncode(t *testing.T) {
	cases := []struct {
		inputLocked  bool
		inputVersion uint64
		output       uint64
	}{
		{false, minVersion, minVersion},
		{false, midVersion, midVersion},
		{false, maxVersion, maxVersion},
		{true, minVersion, 1<<63 | minVersion},
		{true, midVersion, 1<<63 | midVersion},
		{true, maxVersion, 1<<63 | maxVersion},
	}

	for _, tc := range cases {
		encoded := encode(tc.inputLocked, tc.inputVersion)
		if encoded != tc.output {
			t.Errorf("got: %v want: %v", encoded, tc.output)
		}
	}
}

func TestDecode(t *testing.T) {
	cases := []struct {
		locked  bool
		version uint64
	}{
		{false, minVersion},
		{false, midVersion},
		{false, maxVersion},
		{true, minVersion},
		{true, midVersion},
		{true, maxVersion},
	}

	for _, tc := range cases {
		encoded := encode(tc.locked, tc.version)
		locked, version := decode(encoded)
		if locked != tc.locked || version != tc.version {
			t.Errorf("got: (%v, %v) want: (%v, %v)",
				locked, version, tc.locked, tc.version)
		}
	}
}

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
		{false, maxVersion},
		{true, minVersion},
		{true, midVersion},
		{true, maxVersion},
	}

	for _, tc := range cases {
		lock := helperEncode(tc.locked, tc.version)

		ok := lock.tryLock()
		if ok == tc.locked {
			t.Errorf("got: ok = %v, want: ok = %v", ok, !ok)
		}

		locked, version := helperDecode(lock)
		if locked != true || version != tc.version {
			t.Errorf("got: state = (%v, %v) want: state = (%v, %v)",
				locked, version, true, tc.version)
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
		{true, maxVersion},
	}

	for _, tc := range cases {
		lock := helperEncode(tc.locked, tc.version)
		lock.unlock()

		locked, version := helperDecode(lock)
		if locked != false || version != tc.version {
			t.Errorf("got: (%v, %v) want: (%v, %v)",
				locked, version, false, tc.version)
		}
	}
}

func TestUnlockAndUpdate(t *testing.T) {
	cases := []struct {
		locked  bool
		version uint64
		input   uint64
	}{
		{true, minVersion, midVersion},
		{true, midVersion, maxVersion},
	}

	for _, tc := range cases {
		lock := helperEncode(tc.locked, tc.version)
		lock.unlockAndUpdate(tc.input)

		locked, version := helperDecode(lock)
		if locked != false || version != tc.input {
			t.Errorf("got: (%v, %v) want: (%v, %v)",
				locked, version, false, tc.input)
		}
	}
}

func TestSampleLock(t *testing.T) {
	cases := []struct {
		locked  bool
		version uint64
	}{
		{false, minVersion},
		{false, midVersion},
		{false, maxVersion},
		{true, minVersion},
		{true, midVersion},
		{true, maxVersion},
	}

	for _, tc := range cases {
		lock := helperEncode(tc.locked, tc.version)

		locked, version := lock.sampleLock()
		if locked != tc.locked || version != tc.version {
			t.Errorf("got: (%v, %v), want: (%v, %v)",
				locked, version, tc.locked, tc.version)
		}
	}
}

func TestMisuseEncode(t *testing.T) {
	cases := []struct {
		locked  bool
		version uint64
	}{
		{false, 1 << 63},
		{true, 1 << 63},
	}

	for _, tc := range cases {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("Too large version number does not panic.")
			}
		}()

		encode(tc.locked, tc.version)
	}
}

func TestMisuseUnlock(t *testing.T) {
	cases := []struct {
		locked  bool
		version uint64
	}{
		{false, minVersion},
		{false, midVersion},
		{false, maxVersion},
	}

	for _, tc := range cases {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("Unlock of unlocked lock does not panic.")
			}
		}()

		lock := helperEncode(tc.locked, tc.version)
		lock.unlock()
	}
}

func TestMisuseUnlockAndUpdate(t *testing.T) {
	cases := []struct {
		locked  bool
		version uint64
		input   uint64
	}{
		{false, minVersion, midVersion},
		{false, midVersion, maxVersion},
	}

	for _, tc := range cases {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("Unlock of unlocked lock does not panic.")
			}
		}()

		lock := helperEncode(tc.locked, tc.version)
		lock.unlockAndUpdate(tc.input)
	}
}

package stm

import "testing"

const (
	minVersion       = 0
	midVersion       = 1 << 30
	maxVersion       = 1<<63 - 1
	minStateUnlocked = minVersion
	midStateUnlocked = midVersion
	maxStateUnlocked = maxVersion
	minStateLocked   = 1<<63 | minVersion
	midStateLocked   = 1<<63 | midVersion
	maxStateLocked   = 1<<63 | maxVersion
)

func TestEncode(t *testing.T) {
	cases := []struct {
		inputLocked  bool
		inputVersion uint64
		outputState  uint64
	}{
		{false, minVersion, minStateUnlocked},
		{false, midVersion, midStateUnlocked},
		{false, maxVersion, maxStateUnlocked},
		{true, minVersion, minStateLocked},
		{true, midVersion, midStateLocked},
		{true, maxVersion, maxStateLocked},
	}

	for _, tc := range cases {
		encoded := encode(tc.inputLocked, tc.inputVersion)
		if encoded != tc.outputState {
			t.Errorf("got: %v want: %v", encoded, tc.outputState)
		}
	}
}

func TestDecode(t *testing.T) {
	cases := []struct {
		inputLocked  bool
		inputVersion uint64
		// An output pairs should be identical to an input pairs.
	}{
		{false, minVersion},
		{false, midVersion},
		{false, maxVersion},
		{true, minVersion},
		{true, midVersion},
		{true, maxVersion},
	}

	for _, tc := range cases {
		encoded := encode(tc.inputLocked, tc.inputVersion)
		locked, version := decode(encoded)
		if locked != tc.inputLocked || version != tc.inputVersion {
			t.Errorf("got: (%v, %v) want: (%v, %v)",
				locked, version, tc.inputLocked, tc.inputVersion)
		}
	}
}

func TestTryLock(t *testing.T) {
	cases := []struct {
		inputState  uint64
		outputState uint64
		outputOK    bool
	}{
		{minStateUnlocked, minStateLocked, true},
		{midStateUnlocked, midStateLocked, true},
		{maxStateUnlocked, maxStateLocked, true},
		{minStateLocked, minStateLocked, false},
		{midStateLocked, midStateLocked, false},
		{maxStateLocked, maxStateLocked, false},
	}

	for i, tc := range cases {
		lock := (*versionedLock)(&tc.inputState)

		ok := lock.tryLock()
		if ok != tc.outputOK {
			t.Errorf("got: ok = %v, want: ok = %v", ok, tc.outputOK)
		}

		lockState := *(*uint64)(lock)
		if lockState != tc.outputState {
			t.Errorf("%v: got: state = %v want: state = %v",
				i, lockState, tc.outputState)
		}
	}
}

func TestUnlock(t *testing.T) {
	cases := []struct {
		inputState  uint64
		outputState uint64
	}{
		{minStateLocked, minStateUnlocked},
		{midStateLocked, midStateUnlocked},
		{maxStateLocked, maxStateUnlocked},
	}

	for _, tc := range cases {
		lock := (*versionedLock)(&tc.inputState)
		lock.unlock()
		lockState := *(*uint64)(lock)
		if lockState != tc.outputState {
			t.Errorf("got: %v want: %v", lockState, tc.outputState)
		}
	}
}

func TestUnlockAndUpdate(t *testing.T) {
	cases := []struct {
		inputState   uint64
		inputVersion uint64
		// An output should be equal to an inputVersion.
	}{
		{minStateLocked, midVersion},
		{minStateLocked, maxVersion},
		{midStateLocked, minVersion},
		{midStateLocked, maxVersion},
	}

	for _, tc := range cases {
		lock := (*versionedLock)(&tc.inputState)
		lock.unlockAndUpdate(tc.inputVersion)
		lockState := *(*uint64)(lock)
		if lockState != tc.inputState {
			t.Errorf("got: %v want: %v", lockState, tc.inputState)
		}
	}
}

func TestSampleLock(t *testing.T) {
	cases := []struct {
		inputState    uint64
		outputLocked  bool
		outputVersion uint64
	}{
		{minStateUnlocked, false, minVersion},
		{midStateUnlocked, false, midVersion},
		{maxStateUnlocked, false, maxVersion},
		{minStateLocked, true, minVersion},
		{midStateLocked, true, midVersion},
		{maxStateLocked, true, maxVersion},
	}

	for _, tc := range cases {
		lock := (*versionedLock)(&tc.inputState)
		locked, version := lock.sampleLock()
		if locked != tc.outputLocked || version != tc.outputVersion {
			t.Errorf("got: (%v, %v), want: (%v, %v)",
				locked, version, tc.outputLocked, tc.outputVersion)
		}
	}
}

func TestMisuseEncode(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("Too large version number does not panic.")
		}
	}()
	encode(false, 1<<63)
}

func TestMisuseUnlock(t *testing.T) {
	cases := []struct {
		inputState uint64
	}{
		{minStateUnlocked},
		{midStateUnlocked},
		{maxStateUnlocked},
	}

	for _, tc := range cases {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("Unlock of unlocked lock does not panic.")
			}
		}()

		lock := (*versionedLock)(&tc.inputState)
		lock.unlock()
	}
}

func TestMisuseUnlockAndUpdate(t *testing.T) {
	cases := []struct {
		inputState   uint64
		inputVersion uint64
	}{
		{minStateUnlocked, midVersion},
		{midStateUnlocked, maxVersion},
	}

	for _, tc := range cases {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("Unlock of unlocked lock does not panic.")
			}
		}()

		lock := (*versionedLock)(&tc.inputState)
		lock.unlockAndUpdate(tc.inputVersion)
	}
}

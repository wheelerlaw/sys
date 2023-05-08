package unix

import (
	"fmt"
	"testing"
)

func TestBuggyKernel(t *testing.T) {
	var (
		sampleUname Utsname
		isBuggy     bool
		err         error
	)

	testCases := map[string]bool{
		"1.2.3.4":      true,
		"6153":         true,
		"6153.1":       true,
		"6153.141":     true,  // <- and above are < 6153.141.1 and should be reported as buggy
		"6153.141.1":   true,  // Identity case
		"6153.141.1.1": false, // <- and below are > 6153.141.1 and should NOT be reported as buggy
		"6153.141.11":  false,
		"6153.142.1":   false,
		"8080":         false,
	}

	for version, expectedOutcome := range testCases {
		sampleUname = Utsname{}
		copy(sampleUname.Version[:], fmt.Sprintf("hello world xnu-%s~foo/bar", version))
		isBuggy, err = buggyKernel(sampleUname)
		if err != nil {
			t.Fatal(err)
		}
		if isBuggy != expectedOutcome {
			t.Fatalf("Version %s is expected to be buggy: %t, but got %t", version, expectedOutcome, isBuggy)
		}
	}
}

package facts

import (
	"strings"
	"testing"
)

func TestOmittedFor(t *testing.T) {
	for _, major := range []int{16, 17, 18} {
		omitted := omittedFor(major)
		if len(omitted) != 1 || !strings.Contains(omitted[0], "extended statistics") {
			t.Errorf("%d: omitted = %v", major, omitted)
		}
	}
	if omitted := omittedFor(19); len(omitted) != 0 {
		t.Errorf("19: omitted = %v", omitted)
	}
}

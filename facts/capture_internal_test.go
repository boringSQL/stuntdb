package facts

import (
	"bytes"
	"testing"
)

func TestNormaliseRestrict(t *testing.T) {
	first := normaliseRestrict([]byte("\\restrict aaa\nselect 1;\n\\unrestrict aaa\n"))
	second := normaliseRestrict([]byte("\\restrict bbb\nselect 1;\n\\unrestrict bbb\n"))
	if !bytes.Equal(first, second) {
		t.Errorf("two dumps normalise differently:\n%s\n%s", first, second)
	}
	if !bytes.Contains(first, []byte("\\restrict stuntdb\n")) {
		t.Errorf("restrict line = %s", first)
	}
	if !bytes.Contains(first, []byte("\\unrestrict stuntdb\n")) {
		t.Errorf("unrestrict line = %s", first)
	}
	if !bytes.Contains(first, []byte("select 1;\n")) {
		t.Errorf("the dump lost its statements: %s", first)
	}
}

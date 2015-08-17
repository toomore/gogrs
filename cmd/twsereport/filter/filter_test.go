package filter

import "testing"

func TestCheckGroup_String(t *testing.T) {
	for i, v := range AllList {
		t.Log(i, v.No(), v)
	}
}

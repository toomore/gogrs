package utils

import "testing"

func TestHttpCache(*testing.T) {
	hc := NewHttpCache("./.temp", false)
	hc.Get("http://toomore.net/")
}

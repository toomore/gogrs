package utils

import "testing"

func TestHTTPCache(*testing.T) {
	hc := NewHTTPCache("./.temp", false)
	hc.Get("http://toomore.net/?q=%d", true)
}

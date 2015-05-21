package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestHTTPCache(t *testing.T) {
	var dir = filepath.Join(os.TempDir(), fmt.Sprintf("%d", RandInt()))
	defer os.RemoveAll(dir)

	t.Log("TempDir: ", dir)
	hc := NewHTTPCache(dir, false)
	hc.Get("http://toomore.net/?q=%d", true)
	hc.Get("http://toomore.net/?q=%d", true)
}

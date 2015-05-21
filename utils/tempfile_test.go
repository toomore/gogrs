package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestHTTPCache(t *testing.T) {
	var dir = filepath.Join("/Volumes/RamDisk/tmp/.gogrs", fmt.Sprintf("%d", RandInt()))

	t.Log("TempDir: ", dir)
	hc := NewHTTPCache(dir)
	defer os.RemoveAll(hc.Dir)
	hc.Get("http://toomore.net/?q=%d", true)
	hc.Get("http://toomore.net/?q=%d", true)
}

package utils

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func TestHTTPCache(t *testing.T) {
	var dir = filepath.Join("/Volumes/RamDisk/tmp/.gogrs", fmt.Sprintf("%d", RandInt()))

	t.Log("TempDir: ", dir)
	hc := NewHTTPCache(dir, "utf8")
	t.Log("TempDir: ", hc.Dir)
	defer os.RemoveAll(hc.Dir)
	hc.Get("http://toomore.net/?q=%d", true)
	hc.Get("http://toomore.net/?q=%d", true)

	hc.PostForm("http://httpbin.org/post", url.Values{"name": {"Toomore"}})
	hc.PostForm("http://httpbin.org/post", url.Values{"name": {"Toomore"}})

	hccp950 := NewHTTPCache(dir, "cp950")
	hccp950.Get("http://toomore.net/", false)
}

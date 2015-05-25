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

func TestHTTPCache_FlushALL(t *testing.T) {
	hc := NewHTTPCache("/Volumes/RamDisk", "utf8")
	hc.FlushAll()
}

// 目前可以支援 http.Get / http.PostForm 取得資料並儲存
func ExampleHTTPCache() {
	hc := NewHTTPCache("/run/shm/", "utf8") // linux

	ex1, _ := hc.Get("http://httpbin.org/get", false)
	ex2, _ := hc.Get("http://httpbin.org/get?q=%d", true)
	ex3, _ := hc.PostForm("http://httpbin.org/post", url.Values{"name": {"Toomore"}})

	fmt.Printf("%s %s %s", ex1, ex2, ex3)
}

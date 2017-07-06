package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	iconv "github.com/djimenez/iconv-go"
)

// TempFolderName 快取資料夾名稱
const TempFolderName = ".gogrscache"

// HTTPCache net/http 快取功能
type HTTPCache struct {
	Dir            string
	fullpath       string
	iconvConverter func([]byte) []byte
}

// NewHTTPCache New 一個 HTTPCache.
//
// dir 為暫存位置，fromEncoding 來源檔案的編碼，一律轉換為 utf8
func NewHTTPCache(dir string, fromEncoding string) *HTTPCache {
	dir = makeCacheDir(dir)
	return &HTTPCache{
		Dir:            dir,
		fullpath:       filepath.Join(dir, TempFolderName),
		iconvConverter: renderIconvConverter(fromEncoding)}
}

// makeCacheDir 建立快取資料夾
func makeCacheDir(dir string) string {
	var fullpath = filepath.Join(dir, TempFolderName)

	if err := os.Mkdir(fullpath, 0700); os.IsNotExist(err) {
		dir = os.TempDir()
		fullpath = filepath.Join(os.TempDir(), TempFolderName)
		os.Mkdir(fullpath, 0700)
	}
	return dir
}

// Get 透過 http.Get 取得檔案或從暫存中取得檔案
//
// rand 為是否支援網址帶入亂數值，url 需有 '%d' 格式。
func (hc HTTPCache) Get(url string, rand bool) ([]byte, error) {
	filehash := fmt.Sprintf("%x", md5.Sum([]byte(url)))
	var (
		content []byte
		err     error
	)

	if content, err = hc.readFile(filehash); err != nil {
		return hc.saveFile(url, filehash, rand, nil)
	}
	return content, nil
}

// PostForm 透過 http.PostForm 取得檔案或從暫存中取得檔案
func (hc HTTPCache) PostForm(url string, data url.Values) ([]byte, error) {
	hash := md5.New()
	io.WriteString(hash, url)
	io.WriteString(hash, data.Encode())

	var (
		content []byte
		err     error
	)

	filehash := fmt.Sprintf("%x", hash.Sum(nil))
	if content, err = hc.readFile(filehash); err != nil {
		return hc.saveFile(url, filehash, false, data)
	}
	return content, nil
}

// FlushAll 清除快取
func (hc *HTTPCache) FlushAll() {
	os.RemoveAll(hc.fullpath)
	hc.Dir = makeCacheDir(hc.Dir)
}

// readFile 從快取資料裡面取得
func (hc HTTPCache) readFile(filehash string) ([]byte, error) {
	var (
		f   *os.File
		err error
	)
	if f, err = os.Open(filepath.Join(hc.fullpath, filehash)); err == nil {
		defer f.Close()
		return ioutil.ReadAll(f)
	}
	return nil, err
}

// HTTPClient is default func and fixed http too many open files.
var HTTPClient = &http.Client{Transport: &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	Dial: (&net.Dialer{
		Timeout:   0,
		KeepAlive: 0,
	}).Dial,
	TLSHandshakeTimeout: 1 * time.Second,
},
}

// saveFile 從網路取得資料後放入快取資料夾
func (hc HTTPCache) saveFile(url, filehash string, rand bool, data url.Values) ([]byte, error) {
	if rand {
		url = fmt.Sprintf(url, RandInt())
	}

	var (
		content []byte
		err     error
		f       *os.File
		out     []byte
		req     *http.Request
		resp    *http.Response
	)

	if len(data) == 0 {
		// http.Get
		req, err = http.NewRequest("GET", url, nil)
	} else {
		// http.PostForm
		req, err = http.NewRequest("POST", url, strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	if err != nil {
		return out, err
	}

	req.Header.Set("Connection", "close")
	if resp, err = HTTPClient.Do(req); err != nil {
		return out, err
	}
	defer resp.Body.Close()

	if 200 <= resp.StatusCode && resp.StatusCode < 300 {
		if content, err = ioutil.ReadAll(resp.Body); err != nil {
			return out, err
		}

		if f, err = os.Create(filepath.Join(hc.fullpath, filehash)); err != nil {
			return out, err
		}
		defer f.Close()

		out = hc.iconvConverter(content)
		f.Write(out)
	}

	return out, err
}

// renderIconvConverter wrapper function for iconv converter.
func renderIconvConverter(fromEncoding string) func([]byte) []byte {
	if fromEncoding == "utf8" || fromEncoding == "utf-8" {
		return func(str []byte) []byte {
			return str
		}
	}
	return func(content []byte) []byte {
		converter, _ := iconv.NewConverter(fromEncoding, "utf-8")
		var out []byte
		out = make([]byte, len(content)*2)
		_, outLen, _ := converter.Convert(content, out)
		return out[:outLen]
	}
}

// GetOSRamdiskPath try to get RamDisk path.
func GetOSRamdiskPath(goos string) string {
	switch goos {
	case "darwin":
		return "/Volumes/RamDisk/"
	case "linux":
		return "/run/shm/"
	case "":
		return GetOSRamdiskPath(runtime.GOOS)
	default:
		return os.TempDir()
	}
}

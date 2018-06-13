package core

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

type HTTPResponseCache struct {
	sync.RWMutex
	directory string
}

var (
	HTTPResponseCacheInstance *HTTPResponseCache
	once                      sync.Once
)

func (c *HTTPResponseCache) InitCacheStore() error {
	dir, err := ioutil.TempDir("", "chimaera-http-response-cache")
	if err != nil {
		return err
	}
	log.Printf("Saving response caches to: %s\n", dir)
	c.directory = dir
	return nil
}

func (c *HTTPResponseCache) DestroyCacheStore() {
	os.Remove(c.directory)
}

func (c *HTTPResponseCache) WriteResponseCache(r *http.Response) error {
	if !c.cachableRequest(r.Request) {
		return nil
	}
	requestID := c.makeRequestID(r.Request)
	fileName := fmt.Sprintf("%s/%s", c.directory, requestID)
	c.Lock()
	defer c.Unlock()
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	r.Write(f)
	log.Printf("Wrote cache for %s to %s\n", requestID, fileName)
	return nil
}

func (c *HTTPResponseCache) GetResponseCache(r *http.Request) (*http.Response, error) {
	if !c.cachableRequest(r) {
		return nil, nil
	}
	c.Lock()
	defer c.Unlock()
	if !c.hasResponseCache(r) {
		return nil, nil
	}

	requestID := c.makeRequestID(r)
	fileName := fmt.Sprintf("%s/%s", c.directory, requestID)
	_, err := os.Stat(fileName)
	if err != nil && os.IsNotExist(err) {
		return nil, err
	}
	cacheEntry, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer cacheEntry.Close()
	response, err := http.ReadResponse(bufio.NewReader(cacheEntry), r)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *HTTPResponseCache) makeRequestID(r *http.Request) string {
	h := sha1.New()
	io.WriteString(h, r.Method)
	io.WriteString(h, r.URL.String())
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (c *HTTPResponseCache) hasResponseCache(r *http.Request) bool {
	requestID := c.makeRequestID(r)
	fileName := fmt.Sprintf("%s/%s", c.directory, requestID)
	_, err := os.Stat(fileName)
	if err == nil {
		return true
	}
	return false
}

func (c *HTTPResponseCache) cachableRequest(r *http.Request) bool {
	if r.Method == "GET" || r.Method == "HEAD" {
		return true
	}
	return false
}

// http://marcio.io/2015/07/singleton-pattern-in-go/
func GetHTTPResponseCacheInstance() *HTTPResponseCache {
	once.Do(func() {
		HTTPResponseCacheInstance = &HTTPResponseCache{}
		HTTPResponseCacheInstance.InitCacheStore()
	})
	return HTTPResponseCacheInstance
}

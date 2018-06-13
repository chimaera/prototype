package core

import (
	"crypto/tls"
	"github.com/parnurzeal/gorequest"
	"time"
)

func Get(url string, timeout int) (resp gorequest.Response, body string, errs []error) {
	// Make a GET request to the target path
	resp, body, errs = gorequest.New().Proxy("http://localhost:8181").TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Timeout(time.Duration(timeout) * time.Second).Get(url).End()

	return resp, body, errs
}

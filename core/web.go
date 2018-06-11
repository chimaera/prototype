package core

import (
	"crypto/tls"
	"github.com/parnurzeal/gorequest"
	"time"
)

func Get(url string, timeout int) (resp gorequest.Response, body string, errs []error) {
	// Make a GET request to the target path
	resp, body, errs = gorequest.New().TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Timeout(time.Duration(timeout)*time.Second).Get(url).
		// Can be changed later
		Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:31.0) Gecko/20130401 Firefox/31.0").
		End()

	return resp, body, errs
}

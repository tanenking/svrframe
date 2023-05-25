package helper

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/tanenking/svrframe/logx"
)

func DoPost(Domain string, headers map[string]string, data string) (resp *http.Response) {

	req, err := http.NewRequest("POST", Domain, strings.NewReader(data))
	if err != nil {
		logx.ErrorF("DoPost NewRequest err = %v", err)
		return
	}
	client := &http.Client{Timeout: 3 * time.Second}
	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	logx.DebugF("DoPost Domain = %s", Domain)
	logx.DebugF("DoPost data = %s", data)

	resp, err = client.Do(req)
	if err != nil {
		logx.ErrorF("DoPost client.Do err = %v", err)
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		logx.ErrorF("DoPost StatusCode = %d, msg = %s", resp.StatusCode, resp.Status)
		return nil
	}
	return
}
func DoGet(Domain string, headers map[string]string, params map[string]string) (resp *http.Response) {

	Url, err := url.Parse(Domain)
	if err != nil {
		return
	}
	querys := url.Values{}
	for k, v := range params {
		querys.Set(k, v)
	}
	Url.RawQuery = querys.Encode()
	urlPath := Url.String()
	logx.DebugF("urlPath = %s", urlPath)

	req, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		logx.ErrorF("DoGet NewRequest err = %v", err)
		return
	}
	client := &http.Client{Timeout: 3 * time.Second}
	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	resp, err = client.Do(req)
	if err != nil {
		logx.ErrorF("DoGet client.Do err = %v", err)
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		logx.ErrorF("DoGet StatusCode = %d, msg = %s", resp.StatusCode, resp.Status)
		return nil
	}
	return
}

func DoGetDirect(Domain string, headers map[string]string, query string) (resp *http.Response) {

	encode_query := query //url.QueryEscape(query)
	urlPath := Domain
	if len(encode_query) > 0 {
		urlPath += "?" + encode_query
	}
	logx.DebugF("urlPath = %s", urlPath)

	req, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		logx.ErrorF("DoGet NewRequest err = %v", err)
		return
	}
	client := &http.Client{Timeout: 3 * time.Second}
	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	resp, err = client.Do(req)
	if err != nil {
		logx.ErrorF("DoGet client.Do err = %v", err)
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		logx.ErrorF("DoGet StatusCode = %d, msg = %s", resp.StatusCode, resp.Status)
		return nil
	}
	return
}

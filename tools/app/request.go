package app

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
)

func UriRequest(method, uri string, body io.Reader, query url.Values, heads map[string]string) (resp *http.Response, err error) {
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return
	}

	if query != nil {
		u.RawQuery = query.Encode()
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	defer client.CloseIdleConnections()

	return execRequest(client, method, u, body, heads)
}

func execRequest(client *http.Client, method string, u *url.URL, body io.Reader, heads map[string]string) (resp *http.Response, err error) {
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return
	}

	for key, val := range heads {
		reqHeadSet(req, key, val)
	}

	return client.Do(req)
}

func reqHeadSet(req *http.Request, key, value string) {
	req.Header.Set(key, value)
}

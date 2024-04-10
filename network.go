package cmf

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type FetchOptions struct {
	Method     string
	Body       io.Reader
	Headers    map[string]string
	HTTPClient *http.Client
}

var defaultFetchOptions = FetchOptions{
	Method:     http.MethodGet,
	Body:       nil,
	Headers:    nil,
	HTTPClient: http.DefaultClient,
}

func getOpts(options []FetchOptions) FetchOptions {
	opts := defaultFetchOptions
	if len(options) > 0 {
		userOptions := options[0]
		if userOptions.Method != "" {
			opts.Method = userOptions.Method
		}
		if userOptions.Body != nil {
			opts.Body = userOptions.Body
		}
		if userOptions.Headers != nil {
			opts.Headers = userOptions.Headers
		}
		if userOptions.HTTPClient != nil {
			opts.HTTPClient = userOptions.HTTPClient
		}
	}
	return opts
}

/*
Fetch is for making http calls, especially when you need to do a weird on-off
PATCH call or something. Mostly only useful for getting rid of a couple of layers
of `if err != nil` in quick scripts.

TL;DR: Trying to make Fetch happen
*/
func Fetch(url string, options ...FetchOptions) ([]byte, error) {
	opts := getOpts(options)
	req, err := http.NewRequest(opts.Method, url, opts.Body)
	if err != nil {
		return nil, err
	}
	for k, v := range opts.Headers {
		req.Header.Add(k, v)
	}
	res, err := opts.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	// don't report errors on closing connection
	defer func() { _ = res.Body.Close() }()
	return io.ReadAll(res.Body)
}

/*
FetchString is like Fetch, but returns response as a string
*/
func FetchString(url string, options ...FetchOptions) (string, error) {
	bs, err := Fetch(url, options...)
	return string(bs), err
}

/*
FetchJSON is like Fetch, but json decodes the result.
*/
func FetchJSON[T any](url string, options ...FetchOptions) (T, error) {
	var t T
	bs, err := Fetch(url, options...)
	if err != nil {
		return t, err
	}
	err = json.Unmarshal(bs, &t)
	return t, err
}

func PostJSON[T any, U any](url string, body T, options ...FetchOptions) (U, error) {
	var u U

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return u, err
	}
	opts := getOpts(options)
	opts.Body = bytes.NewReader(bodyBytes)
	opts.Method = `POST`
	return FetchJSON[U](url, opts)
}

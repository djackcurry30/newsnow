package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

var defaultClient *http.Client

func init() {
	defaultClient = &http.Client{
		Timeout: 10 * time.Second,
	}
}

type FetchOptions struct {
	Headers   map[string]string
	Body      interface{}
	Timeout   time.Duration
	Retry     int
	Method    string
}

func SetClient(client *http.Client) {
	defaultClient = client
}

func Fetch(url string, options ...FetchOptions) (string, error) {
	opts := FetchOptions{
		Timeout: 10 * time.Second,
		Retry:   3,
		Method:  "GET",
	}
	if len(options) > 0 {
		opts = options[0]
	}

	client := defaultClient
	if opts.Timeout > 0 {
		client = &http.Client{
			Timeout: opts.Timeout,
		}
	}

	var body []byte
	if opts.Body != nil {
		body, _ = json.Marshal(opts.Body)
	}

	req, err := http.NewRequest(opts.Method, url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36")
	for k, v := range opts.Headers {
		req.Header.Set(k, v)
	}

	var lastErr error
	for i := 0; i <= opts.Retry; i++ {
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			bodyBytes, _ := io.ReadAll(resp.Body)
			lastErr = fmt.Errorf("status code %d: %s", resp.StatusCode, string(bodyBytes))
			continue
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = err
			continue
		}
		return string(data), nil
	}

	return "", lastErr
}

func FetchJSON(url string, result interface{}, options ...FetchOptions) error {
	data, err := Fetch(url, options...)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), result)
}

func FetchBytes(url string, options ...FetchOptions) ([]byte, error) {
	data, err := Fetch(url, options...)
	if err != nil {
		return nil, err
	}
	return []byte(data), nil
}

func FindAllStringSubmatch(pattern string, text string) [][]string {
	re := regexp.MustCompile(pattern)
	return re.FindAllStringSubmatch(text, -1)
}

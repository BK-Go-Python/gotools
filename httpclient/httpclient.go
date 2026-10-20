// Package httpclient 提供 Python urllib.request 风格的 HTTP 客户端。
package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// DefaultTimeout 是未指定超时时的默认值。
const DefaultTimeout = 30 * time.Second

// Options 是请求选项。
type Options struct {
	Timeout time.Duration
	Headers map[string]string
	Query   url.Values
}

// Response 是 HTTP 响应。
type Response struct {
	StatusCode int
	Status     string
	Headers    http.Header
	Body       []byte
	URL        string
}

// OK 判断状态码是否为 2xx。
func (r *Response) OK() bool {
	return r != nil && r.StatusCode >= 200 && r.StatusCode < 300
}

// Text 返回响应体字符串。
func (r *Response) Text() string {
	if r == nil {
		return ""
	}
	return string(r.Body)
}

// JSON 把响应体解析为 v。
func (r *Response) JSON(v any) error {
	if r == nil {
		return fmt.Errorf("httpclient.Response: nil response")
	}
	if err := json.Unmarshal(r.Body, v); err != nil {
		return fmt.Errorf("httpclient.Response.JSON: %w", err)
	}
	return nil
}

// Get 发送 GET 请求。
func Get(rawURL string, opts ...Options) (*Response, error) {
	return Request(http.MethodGet, rawURL, nil, opts...)
}

// Post 发送 POST 请求。
func Post(rawURL string, body []byte, opts ...Options) (*Response, error) {
	return Request(http.MethodPost, rawURL, bytes.NewReader(body), opts...)
}

// PostJSON 把 value 编码为 JSON 并发送 POST 请求。
func PostJSON(rawURL string, value any, opts ...Options) (*Response, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("httpclient.PostJSON: %w", err)
	}
	option := firstOptions(opts)
	if option.Headers == nil {
		option.Headers = make(map[string]string)
	}
	if _, ok := option.Headers["Content-Type"]; !ok {
		option.Headers["Content-Type"] = "application/json"
	}
	return Request(http.MethodPost, rawURL, bytes.NewReader(data), option)
}

// GetJSON 发送 GET 请求并把响应 JSON 解析到 out。
func GetJSON(rawURL string, out any, opts ...Options) error {
	response, err := Get(rawURL, opts...)
	if err != nil {
		return err
	}
	return response.JSON(out)
}

// Request 发送任意方法请求。
func Request(method, rawURL string, body io.Reader, opts ...Options) (*Response, error) {
	option := firstOptions(opts)
	target, err := applyQuery(rawURL, option.Query)
	if err != nil {
		return nil, fmt.Errorf("httpclient.Request: %w", err)
	}
	request, err := http.NewRequest(method, target, body)
	if err != nil {
		return nil, fmt.Errorf("httpclient.Request: %w", err)
	}
	for key, value := range option.Headers {
		request.Header.Set(key, value)
	}
	timeout := option.Timeout
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	client := &http.Client{Timeout: timeout}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("httpclient.Request: %w", err)
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("httpclient.Request: %w", err)
	}
	return &Response{
		StatusCode: response.StatusCode,
		Status:     response.Status,
		Headers:    response.Header,
		Body:       data,
		URL:        response.Request.URL.String(),
	}, nil
}

func firstOptions(opts []Options) Options {
	if len(opts) > 0 {
		return opts[0]
	}
	return Options{}
}

func applyQuery(rawURL string, query url.Values) (string, error) {
	if len(query) == 0 {
		return rawURL, nil
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	existing := parsed.Query()
	for key, values := range query {
		existing[key] = append(existing[key], values...)
	}
	parsed.RawQuery = existing.Encode()
	return parsed.String(), nil
}

package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-framework/util/helper"
	"go-framework/util/tracer"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client http.Client
type Client struct {
	domain string
	client *http.Client
}

// Options 客户端函数选项
type Options func(*Option)

// Option 客户端选项
type Option struct {
	timeout   time.Duration
	transport http.RoundTripper
	headers   map[string]string
}

// WithTimeout 设置HTTP客户端的超时
func WithTimeout(timeout time.Duration) Options {
	return func(o *Option) {
		o.timeout = timeout
	}
}

// WithTransport 设置HTTP客户端的传输
func WithTransport(transport http.RoundTripper) Options {
	return func(o *Option) {
		o.transport = transport
	}
}

// WithHeaders 设置HTTP请求头
func WithHeaders(headers map[string]string) Options {
	return func(o *Option) {
		o.headers = headers
	}
}

// DataForm 设置表单
func DataForm(data map[string]string) []byte {
	list := url.Values{}
	for k, v := range data {
		list.Add(k, v)
	}
	return []byte(list.Encode())
}

// DataJSON 设置JSON
func DataJSON(data interface{}) []byte {
	if data == nil {
		return nil
	}
	postJSON, _ := helper.Marshal(data)
	return postJSON
}

// NewClient 新的客户端实例
func NewClient(domain string) *Client {
	return &Client{
		domain: domain,
		client: &http.Client{},
	}
}

func applyOptions(options []Options) *Option {
	o := &Option{
		transport: http.DefaultTransport,
		headers:   make(map[string]string),
	}
	for _, option := range options {
		option(o)
	}
	return o
}

// Get 执行GET请求
func (c *Client) Get(ctx context.Context, endpoint string, params map[string]string, options ...Options) ([]byte, error) {
	reqURL, err := url.Parse(c.domain + endpoint)
	if err != nil {
		return nil, err
	}

	o := applyOptions(options)

	query := reqURL.Query()
	for key, value := range params {
		query.Set(key, value)
	}
	reqURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, err
	}

	headers := tracer.MapCarrier(ctx)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	for key, value := range o.headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{
		Timeout:   o.timeout,
		Transport: o.transport,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET请求失败，状态代码为: %d, msg: %s", resp.StatusCode, body)
	}

	return body, nil
}

// Post 执行POST请求
func (c *Client) Post(ctx context.Context, endpoint string, data []byte, options ...Options) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.domain+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	headers := tracer.MapCarrier(ctx)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	req.Header.Set("Content-Type", "application/json")

	o := applyOptions(options)
	for key, value := range o.headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{
		Timeout:   o.timeout,
		Transport: o.transport,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("POST请求失败，状态代码为: %d, msg: %s", resp.StatusCode, body)
	}
	return body, nil
}

// Put 执行PUT请求
func (c *Client) Put(ctx context.Context, endpoint string, data interface{}, options ...Options) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	o := applyOptions(options)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.domain+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	headers := tracer.MapCarrier(ctx)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range o.headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{
		Timeout:   o.timeout,
		Transport: o.transport,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("PUT请求失败，状态代码为: %d, msg: %s", resp.StatusCode, body)
	}

	return body, nil
}

// Delete 执行DELETE请求
func (c *Client) Delete(ctx context.Context, endpoint string, params map[string]string, options ...Options) ([]byte, error) {
	reqURL, err := url.Parse(c.domain + endpoint)
	if err != nil {
		return nil, err
	}

	query := reqURL.Query()
	for key, value := range params {
		query.Set(key, value)
	}
	reqURL.RawQuery = query.Encode()

	o := applyOptions(options)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, reqURL.String(), nil)
	if err != nil {
		return nil, err
	}
	headers := tracer.MapCarrier(ctx)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	for key, value := range o.headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{
		Timeout:   o.timeout,
		Transport: o.transport,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DELETE请求失败，状态代码为: %d, msg: %s", resp.StatusCode, body)
	}

	return body, nil
}

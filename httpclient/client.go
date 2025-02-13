package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	token      string
}

func NewClient(baseURL, token string, timeout *time.Duration) *Client {
	var t time.Duration
	if timeout != nil {
		t = *timeout
	} else {
		t = 30 * time.Second
	}
	return &Client{
		httpClient: &http.Client{
			Timeout: t,
		},
		baseURL: baseURL,
		token:   token,
	}
}

func (c *Client) call(ctx context.Context, apiPath, method string, queryParams url.Values, requestBody, res interface{}, customHeaders map[string]string) error {
	var body io.Reader

	if requestBody != nil {
		if reader, ok := requestBody.(io.Reader); ok {
			body = reader
		} else {
			jsonParams, err := json.Marshal(requestBody)
			if err != nil {
				return err
			}
			body = bytes.NewBuffer(jsonParams)
		}
	}

	if customHeaders == nil {
		customHeaders = map[string]string{"Content-Type": "application/json"}
	}

	req, err := c.newRequest(ctx, apiPath, method, customHeaders, queryParams, body)
	if err != nil {
		return err
	}
	return c.do(req, res)
}

func (c *Client) newRequest(ctx context.Context, apiPath string, method string, headers map[string]string, queryParams url.Values, body io.Reader) (*http.Request, error) {
	var u *url.URL
	var err error

	if strings.HasPrefix(apiPath, "https://") {
		u, err = url.Parse(apiPath)
		if err != nil {
			return nil, err
		}
	} else {
		base, err := url.Parse(c.baseURL)
		if err != nil {
			return nil, err
		}
		if !strings.HasSuffix(base.Path, "/") {
			base.Path += "/"
		}
		u = base.ResolveReference(&url.URL{Path: apiPath})
	}

	if queryParams != nil {
		existingQuery := u.Query()
		for key, values := range queryParams {
			for _, value := range values {
				existingQuery.Add(key, value)
			}
		}
		u.RawQuery = existingQuery.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

func (c *Client) do(req *http.Request, res interface{}) error {
	response, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if res == nil {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}
		fmt.Printf("response body: %s\n", string(bodyBytes))

		return nil
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}
		return fmt.Errorf("HTTP %d: %s", response.StatusCode, string(bodyBytes))
	}

	return json.NewDecoder(response.Body).Decode(&res)
}

package payjp

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
	
	"github.com/payjp/payjp-cli/internal/version"
)

// Client is the API client used to sent requests to Payjp's API.
type Client struct {
	baseURL    *url.URL
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Client with the given API key.
func NewClient(baseURLStr, apiKey string) (*Client, error) {
	baseURL, err := url.Parse(baseURLStr)
	if err != nil {
		return nil, err
	}

	c := &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
	}

	return c, nil
}

func (c *Client) client() *http.Client {
	if c.httpClient == nil {
		c.httpClient = newHTTPClient()
	}
	return c.httpClient
}

func (c *Client) PerformRequest(ctx context.Context, method, path string, params url.Values) (*http.Response, error) {
	url, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	url = c.baseURL.ResolveReference(url)

	if method == http.MethodGet {
		query := url.Query()
		for k, v := range params {
			query.Set(k, v[0])
		}
		url.RawQuery = query.Encode()
	}

	req, err := http.NewRequest(method, url.String(), strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "payjp-cli")
	req.Header.Set("X-Client-User-Agent", version.GetUserAgent())

	if c.apiKey != "" {
		req.SetBasicAuth(c.apiKey, "")
	}

	if ctx != nil {
		req = req.WithContext(ctx)
	}

	resp, err := c.client().Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func newHTTPClient() *http.Client {
	httpTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return &http.Client{
		Transport: httpTransport,
	}
}

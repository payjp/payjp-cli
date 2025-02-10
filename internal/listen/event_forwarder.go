package listen

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	pb "github.com/payjp/payjp-cli/gen/proto"
)

type EventForwarder struct {
	forwardURL *url.URL
	httpClient *http.Client
}

func NewEventForwarder(forwardTo string, skipVerify bool) (*EventForwarder, error) {
	forwardURL, err := parseForwardURL(forwardTo)
	if err != nil {
		return nil, err
	}
	return &EventForwarder{
		forwardURL: forwardURL,
		httpClient: newHTTPClient(skipVerify),
	}, nil
}

func (f *EventForwarder) ForwardEvent(res *pb.PayjpEventResponse) error {
	if f.forwardURL == nil {
		return nil
	}

	body, err := json.Marshal(res.PayjpEvent)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, f.forwardURL.String(), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	for _, header := range res.Headers {
		if strings.ToLower(header.Key) == "host" {
			req.Host = header.Value
		} else {
			req.Header.Add(header.Key, header.Value)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func parseForwardURL(forwardTo string) (*url.URL, error) {
	if forwardTo == "" {
		return nil, nil
	}

	if strings.HasPrefix(forwardTo, ":") {
		forwardTo = "http://localhost" + forwardTo
	} else if strings.Trim(forwardTo, "0123456789") == "" {
		forwardTo = "http://localhost" + ":" + forwardTo
	} else if !strings.HasPrefix(forwardTo, "http://") && !strings.HasPrefix(forwardTo, "https://") {
		forwardTo = "http://" + forwardTo
	}

	u, err := url.Parse(forwardTo)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func newHTTPClient(skipVerify bool) *http.Client {
	httpTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipVerify,
		},
	}

	return &http.Client{
		Transport: httpTransport,
	}
}

package login

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/payjp/payjp-cli/internal/payjp"
)

type PayjpCliAuthResponse struct {
	BrowserURL  string `json:"browser_url"`
	PollURL     string `json:"poll_url"`
	PairingCode string `json:"pairing_code"`
}

type PayjpCliAuthPollingResponse struct {
	Token              string `json:"token"`
	PairingCode        string `json:"pariing_code"`
	Redeemed           bool   `json:"redeemed"`
	AccountID          string `json:"account_id"`
	AccountDisplayName string `json:"account_display_name"`
	TestModeSecretKey  string `json:"test_mode_secret_key"`
}

type AuthResult struct {
	AccountID          string
	AccountDisplayName string
	TestModeSecretKey  string
	Err                error
}

// CallAuth calls the /payjpcli/auth endpoint to start the authentication process
func CallAuth(ctx context.Context, client *payjp.Client) (*PayjpCliAuthResponse, error) {
	res, err := client.PerformRequest(ctx, "POST", "/payjpcli/auth", url.Values{})
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected http status code: %d %s", res.StatusCode, string(bodyBytes))
	}

	var response *PayjpCliAuthResponse
	jsonErr := json.Unmarshal(bodyBytes, &response)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return response, nil
}

// PollingAuthResult polls the pollingURL endpoint until the authentication process is confirmed
func PollingAuthResult(ctx context.Context, client *payjp.Client, pollURL string, pollingCh chan<- *AuthResult) {
	MAX_POLLING_COUNT := 60
	count := 0

	for {
		res, err := client.PerformRequest(ctx, "GET", pollURL, url.Values{})
		if err != nil {
			pollingCh <- &AuthResult{
				Err: err,
			}
			return
		}

		defer res.Body.Close()
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			pollingCh <- &AuthResult{
				Err: err,
			}
			return
		}

		if res.StatusCode != http.StatusOK {
			pollingCh <- &AuthResult{
				Err: fmt.Errorf("unexpected http status code: %d %s", res.StatusCode, string(bodyBytes)),
			}
			return
		}

		var response *PayjpCliAuthPollingResponse
		jsonErr := json.Unmarshal(bodyBytes, &response)
		if jsonErr != nil {
			pollingCh <- &AuthResult{
				Err: jsonErr,
			}
			return
		}

		if response.Redeemed {
			pollingResult := &AuthResult{
				AccountID:          response.AccountID,
				AccountDisplayName: response.AccountDisplayName,
				TestModeSecretKey:  response.TestModeSecretKey,
			}
			pollingCh <- pollingResult
			return
		}

		count++
		if count > MAX_POLLING_COUNT {
			pollingCh <- &AuthResult{
				Err: fmt.Errorf("timed out. Please try again."),
			}
			return
		}

		time.Sleep(2 * time.Second)
	}
}

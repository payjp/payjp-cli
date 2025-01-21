package listen

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/payjp/payjp-cli/internal/payjp"
)

type PayjpCliSessionStatus string

const (
	PayjpCliSessionStatusUnused    PayjpCliSessionStatus = "unused"
	PayjpCliSessionStatusActivated PayjpCliSessionStatus = "activated"
	PayjpCliSessionStatusDisposed  PayjpCliSessionStatus = "disposed"
)

type PayjpCliSession struct {
	ID      string
	Status  PayjpCliSessionStatus
	Expires int64
}

func CreateCliSession(ctx context.Context, client *payjp.Client) (*PayjpCliSession, error) {
	res, err := client.PerformRequest(ctx, "POST", "/v1/payjpcli/sessions", "")
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

	var session *PayjpCliSession
	jsonErr := json.Unmarshal(bodyBytes, &session)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return session, nil
}

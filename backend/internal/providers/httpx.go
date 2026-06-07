package providers

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// httpDoer is the minimal HTTP surface the diligence/strategy clients need
// (satisfied by *http.Client; trivially mockable in tests).
type httpDoer interface {
	Do(*http.Request) (*http.Response, error)
}

func newRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, method, url, body)
}

// doJSON executes the request and returns the body, treating >=400 as an error.
func doJSON(c httpDoer, req *http.Request) ([]byte, error) {
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}

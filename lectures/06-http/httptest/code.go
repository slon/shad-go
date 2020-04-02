package httptest

import (
	"context"
	"io/ioutil"
	"net/http"
	"strconv"
)

func NewAPICLient(baseURL string) *APIClient {
	if baseURL == "" {
		baseURL = BaseURLProd
	}
	return &APIClient{baseURL: baseURL, httpc: new(http.Client)}
}

const (
	BaseURLProd = "https://github.com/api"
)

type APIClient struct {
	baseURL string
	httpc   *http.Client
}

func (c *APIClient) GetReposCount(ctx context.Context, userID string) (int, error) {
	url := c.baseURL + "/users/" + userID + "/repos/count"
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := c.httpc.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(string(body))
}

// OMIT

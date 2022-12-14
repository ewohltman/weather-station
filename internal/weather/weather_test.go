package weather_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

type mockHTTPClient struct {
	Response []byte
}

func (httpClient *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{Body: io.NopCloser(bytes.NewReader(httpClient.Response))}, nil
}

func TestResponse_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	// TODO
}

func TestResponseError_Error(t *testing.T) {
	t.Parallel()

	// TODO
}

func TestNewAPIClient(t *testing.T) {
	t.Parallel()

	// TODO
}

func TestAPIClient_Query(t *testing.T) {
	t.Parallel()

	// TODO
}

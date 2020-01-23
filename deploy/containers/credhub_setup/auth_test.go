package main

import (
	"bytes"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMakeHTTPClientWithCA(t *testing.T) {
	t.Parallel()
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	}))
	defer server.Close()

	serverURL, err := url.Parse(server.URL)
	require.NoError(t, err, "error parsing server url")

	certBytes := bytes.Buffer{}
	pem.Encode(&certBytes, &pem.Block{
		Type: "CERTIFICATE",
		Bytes: server.Certificate().Raw,
	})
	client, err := makeHTTPClientWithCA(serverURL.Hostname(), certBytes.Bytes())
	require.NoError(t, err, "failed to make HTTP client")
	
	resp, err := client.Get(server.URL)
	require.NoError(t, err, "error fetching from test server")
	require.GreaterOrEqual(t, resp.StatusCode, 200, "unexpected status: %s", resp.Status)
	require.Less(t, resp.StatusCode, 300, "unexpected status: %s", resp.Status)
}

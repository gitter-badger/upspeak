package matrix

import (
	"net/http"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
)

const (
	baseURL = "https://matrix.org"
	userID  = "@test:upspeak.net"
	token   = "i19291k2imimimsidxmdi1m9k29oi3kmimimaisds=="
)

func apiPath(paths ...string) string {
	return strings.Join([]string{baseURL, prefixPath, strings.Join(paths, "/")}, "/")
}

func TestNewClient(t *testing.T) {
	c, err := newClient(baseURL, userID)
	if err != nil {
		t.Errorf("Error creating client. Err: %s", err.Error())
	}
	if c.baseURL.Path != prefixPath {
		t.Errorf("Client Prefix path not set correctly")
	}
	c.token(token)
	if c.accessToken != token {
		t.Errorf("Client Access token not set correctly")
	}
	// c.httpClient
}

type testRes struct {
	Message string `json:"message"`
}

func TestSend(t *testing.T) {
	c, err := newClient(baseURL, userID)
	if err != nil {
		t.Errorf("Error creating client. Err: %s", err.Error())
	}
	httpmock.ActivateNonDefault(c.httpClient.GetClient())
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", apiPath("test"),
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, testRes{Message: "test123"})
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	var tres testRes
	err = c.send("GET", "test", nil, nil, &tres)
	if err != nil {
		t.Error("Error sending request")
	}
	if tres.Message != "test123" {
		t.Error("client.send() returned invalid response")
	}
}

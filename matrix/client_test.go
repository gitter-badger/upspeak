package matrix

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestNewClient(t *testing.T) {
	const (
		baseURL = "https://matrix.org"
		userID  = "@test:upspeak.net"
		token   = "sidji1jwimdi939293dmaimdi1m3imiemamsqoef111wma=="
	)
	c, err := NewClient(baseURL, userID)
	if err != nil {
		t.Errorf("Error creating client. Err: %s", err.Error())
	}
	if c.baseURL.Path != "/_matrix/client/r0" {
		t.Errorf("Client Prefix path not set correctly")
	}
	c.Token(token)
	if c.accessToken != token {
		t.Errorf("Client Access token not set correctly")
	}
	if c.UserID() != userID {
		t.Errorf("Invalid user ID returned by ")
	}
}

type testRes struct {
	Message string `json:"message"`
}

func TestAPIPath(t *testing.T) {
	c, err := NewClient("https://example.com", "@test:example.com")
	if err != nil {
		t.Error("Error creating new client")
	}
	expectedPath := "https://example.com/_matrix/client/r0/test"
	testPath := c.apiPath("test")
	if testPath != expectedPath {
		t.Errorf("Invalid API Path")
	}
}

func TestSend(t *testing.T) {
	var (
		baseURL = "https://matrix.org"
		userID  = "@test:upspeak.net"
	)
	c, err := NewClient(baseURL, userID)
	if err != nil {
		t.Errorf("Error creating client. Err: %s", err.Error())
	}
	httpmock.ActivateNonDefault(c.httpClient.GetClient())
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", c.apiPath("test"),
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, testRes{Message: "test123"})
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	var tres testRes
	err = c.send(request{
		method:  "GET",
		subPath: "test",
	}, &tres)
	if err != nil {
		t.Error("Error sending request")
	}
	if tres.Message != "test123" {
		t.Error("client.send() returned invalid response")
	}
}

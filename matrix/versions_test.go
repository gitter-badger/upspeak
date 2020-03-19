package matrix

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestVersionSuccess(t *testing.T) {
	var (
		baseURL     = "https://matrix.org"
		userID      = "@test:upspeak.net"
		expectedRes = VersionsResponse{
			Versions: []string{"r0.0.1"},
			UnstableFeatures: map[string]bool{
				"net.upspeak.testing": true,
			},
		}
	)

	c, err := NewClient(baseURL, userID)
	if err != nil {
		t.Errorf("Error creating client. Err: %s", err.Error())
	}

	httpmock.ActivateNonDefault(c.httpClient.GetClient())
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", c.apiPath("versions"),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, expectedRes)
		},
	)

	res, err := c.Versions()

	if err != nil {
		t.Error("Error from versions", err)
	}
	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Invalid response object received using Client.Versions()")
	}
}

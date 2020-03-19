package matrix

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestLoginPasswordSuccess(t *testing.T) {
	var (
		baseURL     = "https://matrix.org"
		userID      = "@test:upspeak.net"
		expectedRes = LoginPasswordResponse{
			UserID:      userID,
			AccessToken: "1wimdima9931mdimaiwjdi13k9dm3==",
			DeviceID:    "Upspeak",
			WellKnown: DiscoveryInformation{
				HomeServer:     ServerInformation{BaseURL: baseURL},
				IdentityServer: ServerInformation{BaseURL: "https://id.example.com"},
			},
		}
	)

	c, err := NewClient(baseURL, userID)
	if err != nil {
		t.Errorf("Error creating client. Err: %s", err.Error())
	}

	httpmock.ActivateNonDefault(c.httpClient.GetClient())
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", c.apiPath("login"),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, expectedRes)
		},
	)

	res, err := c.LoginPassword("goodpassword")

	if err != nil {
		t.Error("Error from login", err)
	}
	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Invalid response object received using Client.LoginPassword()")
	}
}

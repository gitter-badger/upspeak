package matrix

import (
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

// client is the main type for handling requests to the Matrix client-server API
type client struct {
	baseURL     *url.URL
	userID      string
	accessToken string
	httpClient  *resty.Client
}

// token adds an accessToken to Client
func (c *client) token(t string) *client {
	c.accessToken = t
	return c
}

// send prepares and sends a HTTP request
func (c *client) send(method string, subPath string, params url.Values, reqBody interface{}, resBody interface{}) error {
	u, err := url.Parse(c.baseURL.String())
	if err != nil {
		return nil
	}
	u.Path = path.Join(u.Path, subPath)
	r := c.httpClient.R().SetContentLength(true)
	r.URL = u.String()
	if c.accessToken != "" {
		r.SetAuthToken(c.accessToken)
	}
	if reqBody != nil {
		r.SetBody(reqBody)
	}
	if params != nil {
		r.SetQueryParamsFromValues(params)
	}
	r.Method = method
	if resBody != nil {
		r.SetResult(resBody)
	}
	res, err := r.Send()
	if err != nil && res.IsError() {
		return err
	}
	return nil
}

func (c *client) apiPath(paths ...string) string {
	b := c.baseURL.String()
	b = b + "/" + strings.Trim(path.Join(paths...), "/")
	return b
}

// newClient creates a new HTTP client to send requests to the Matrix client-server API
func newClient(hsURL, userID string) (*client, error) {
	u, err := url.Parse(hsURL)
	if err != nil {
		return nil, err
	}
	u.Path = "/_matrix/client/r0"
	c := client{
		baseURL:    u,
		userID:     userID,
		httpClient: resty.New().SetTimeout(15 * time.Second),
	}
	return &c, nil
}

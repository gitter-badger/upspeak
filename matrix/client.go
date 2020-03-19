package matrix

import (
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client is the main type for handling requests to the Matrix client-server API
type Client struct {
	baseURL     *url.URL
	userID      string
	accessToken string
	httpClient  *resty.Client
}

// Token adds an accessToken to Client
func (c *Client) Token(t string) *Client {
	c.accessToken = t
	return c
}

// UserID returns the user ID associated with the current client
func (c *Client) UserID() string {
	return c.userID
}

// Internal request structure used in call to Client.send()
type request struct {
	method  string
	subPath string
	params  url.Values
	body    interface{}
}

// send prepares and sends a HTTP request
func (c *Client) send(req request, resBody interface{}) error {
	u, err := url.Parse(c.baseURL.String())
	if err != nil {
		return nil
	}
	u.Path = path.Join(u.Path, req.subPath)
	r := c.httpClient.R().SetContentLength(true)
	r.URL = u.String()
	if c.accessToken != "" {
		r.SetAuthToken(c.accessToken)
	}
	if req.body != nil {
		r.SetBody(req.body)
	}
	if req.params != nil {
		r.SetQueryParamsFromValues(req.params)
	}
	r.Method = req.method
	if resBody != nil {
		r.SetResult(resBody)
	}
	res, err := r.Send()
	if err != nil && res.IsError() {
		return err
	}
	return nil
}

func (c *Client) apiPath(paths ...string) string {
	b := c.baseURL.String()
	b = b + "/" + strings.Trim(path.Join(paths...), "/")
	return b
}

// NewClient creates a new HTTP client to send requests to the Matrix client-server API
func NewClient(hsURL, userID string) (*Client, error) {
	u, err := url.Parse(hsURL)
	if err != nil {
		return nil, err
	}
	u.Path = "/_matrix/client/r0"
	c := &Client{
		baseURL:    u,
		userID:     userID,
		httpClient: resty.New().SetTimeout(15 * time.Second),
	}
	return c, nil
}

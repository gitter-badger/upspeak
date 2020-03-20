package matrix

import (
	"net/url"
	"path"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client is the main type for handling requests to the Matrix client-server API
type Client struct {
	host        string
	prefixPath  string
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
	method     string
	subPath    string
	prefixPath string
	params     url.Values
	body       interface{}
}

// send prepares and sends a HTTP request
func (c *Client) send(req request, resBody interface{}) error {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
	}
	if req.prefixPath != "" {
		u.Path = path.Join(req.prefixPath, req.subPath)
	} else {
		u.Path = path.Join(c.prefixPath, req.subPath)
	}
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

// NewClient creates a new HTTP client to send requests to the Matrix client-server API
func NewClient(hsURL, userID string) (*Client, error) {
	u, err := url.Parse(hsURL)
	if err != nil {
		return nil, err
	}
	c := &Client{
		host:       u.Host,
		prefixPath: "/_matrix/client/r0",
		userID:     userID,
		httpClient: resty.New().SetTimeout(15 * time.Second),
	}
	return c, nil
}

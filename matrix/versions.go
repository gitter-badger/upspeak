package matrix

// VersionsResponse holds the structure for call to `/versions` in the Matrix Client-Server API
type VersionsResponse struct {
	Versions         []string        `json:"versions"`
	UnstableFeatures map[string]bool `json:"UnstableFeatures"`
}

// Versions returns versions of the Matrix specification supported by the Matrix Server
func (c *Client) Versions() (VersionsResponse, error) {
	var res VersionsResponse
	req := request{
		method:  "GET",
		subPath: "versions",
	}
	err := c.send(req, &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

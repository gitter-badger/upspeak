package matrix

// loginPasswordRequest is the structure used to send a login password request
type loginPasswordRequest struct {
	Type       string         `json:"type"`
	Identifier UserIdentifier `json:"identifier"`
	Password   string         `json:"password"`
	DeviceID   string         `json:"device_id"`
	DeviceName string         `json:"initial_device_display_name"`
}

// LoginPasswordResponse is the structure of response for login password
// requests
type LoginPasswordResponse struct {
	UserID      string               `json:"user_id"`
	AccessToken string               `json:"access_token"`
	DeviceID    string               `json:"device_id"`
	WellKnown   DiscoveryInformation `json:"well_known,omitempty"`
}

// LoginPassword sends a login request to Matrix with a password
func (c *Client) LoginPassword(password string) (LoginPasswordResponse, error) {
	var res LoginPasswordResponse
	req := request{
		method:  "POST",
		subPath: "login",
		body: loginPasswordRequest{
			Type: "m.login.password",
			Identifier: UserIdentifier{
				Type: "m.id.user",
				User: c.userID,
			},
			Password:   password,
			DeviceName: "Upspeak",
		},
	}
	err := c.send(req, &res)
	if err != nil {
		return res, err
	}
	c.Token(res.AccessToken)
	return res, nil
}

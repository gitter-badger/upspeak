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
func LoginPassword(homeserverURL, userID, password string) (LoginPasswordResponse, error) {
	var res LoginPasswordResponse
	c, err := newClient(homeserverURL, userID)
	if err != nil {
		return res, err
	}
	req := loginPasswordRequest{
		Type: "m.login.password",
		Identifier: UserIdentifier{
			Type: "m.id.user",
			User: userID,
		},
		Password:   password,
		DeviceName: "Upspeak",
	}
	if err = c.send("POST", "login", nil, req, res); err != nil {
		return res, err
	}
	return res, nil
}

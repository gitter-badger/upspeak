package matrix

type UserIdentifier struct {
	Type string `json:type`
	User string `json:user`
}

type ServerInformation struct {
	BaseURL string `json:base_url`
}

type DiscoveryInformation struct {
	HomeServer     ServerInformation `json:m.homeserver`
	IdentityServer ServerInformation `json:m.identity_server`
}

type LoginPasswordRequest struct {
	Type                     string         `json:type`
	Identifier               UserIdentifier `json:identifier`
	Password                 string         `json:password`
	DeviceID                 string         `json:device_id`
	InitialDeviceDisplayName string         `json:initial_device_display_name`
}

type LoginPasswordResponse struct {
	UserID      string               `json:user_id`
	AccessToken string               `json:access_token`
	DeviceID    string               `json:device_id`
	WellKnown   DiscoveryInformation `json:well_known,omitempty`
}

func LoginPassword(homeserverURL, userID, password string) (LoginPasswordResponse, error) {
	var res LoginPasswordResponse
	c, err := newClient(homeserverURL, userID)
	if err != nil {
		return res, err
	}
	req := LoginPasswordRequest{
		Type: "m.login.password",
		Identifier: UserIdentifier{
			Type: "m.id.user",
			User: userID,
		},
		Password:                 password,
		InitialDeviceDisplayName: "Upspeak",
	}
	if err = c.send("POST", "login", nil, req, res); err != nil {
		return res, err
	}
	return res, err
}

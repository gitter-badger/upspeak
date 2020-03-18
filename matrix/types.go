package matrix

// UserIdentifier is the type to hold a User Identifier object in Matrix
type UserIdentifier struct {
	Type string `json:"type"`
	User string `json:"user"`
}

// ServerInformation holds the different types of Server Information data types
// of Matrix
type ServerInformation struct {
	BaseURL string `json:"base_url"`
}

// DiscoveryInformation holds the Discovery Information data type of Matrix
type DiscoveryInformation struct {
	HomeServer     ServerInformation `json:"m.homeserver"`
	IdentityServer ServerInformation `json:"m.identity_server"`
}

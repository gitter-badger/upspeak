package matrix

import (
	"strings"
)

// ID represents a Parsed ID
type ID struct {
	Username string `json:"username"`
	Hostname string `json:"hostname"`
}

func (id *ID) String() string {
	return id.Username + ":" + id.Hostname
}

// ParseID parses a Matrix user ID string into id and hostname
func ParseID(userID string) ID {
	s := strings.Split(userID, ":")
	return ID{
		Username: s[0],
		Hostname: s[1],
	}
}

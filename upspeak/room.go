package upspeak

// Room holds the data structure for a Room in Upspeak. It maps to a room in Matrix
type Room struct {
	ID    string  `json:"id"`
	Nodes []*Node // Holds an array of Nodes that belong to this room
}

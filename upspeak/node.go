package upspeak

import (
	"github.com/upspeak/upspeak/matrix"
)

// Content holds the parsed JSON content for Node
type Content interface{}

// Node is the unit of information in Upspeak
type Node struct {
	ID      string    `json:"id"`
	Content Content   `json:"content"`
	Type    string    `json:"type"`
	Sender  matrix.ID `json:"sender"`

	// Relations
	PreviousNode *Node // Relation to the node before this one
	SourceNode   *Node // The origin of the thread
	Room         *Room // The room this Node belongs to
}

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

	// Graph
	PreviousNode  *Node   // Relation to the node before this one
	SourceNode    *Node   // The origin of the thread
	AdjacentNodes []*Node // Nodes that have been forked from this node

	// Part of
	Room *Room // The room this Node belongs to
}

// NewNode creates a new Node instance. It does not persist the created Node.
// Call Node.Save() on the returned Node to persist it.
func NewNode(id string, content Content) Node {
	return Node{
		ID:      id,
		Content: content,
	}
}

// Save persists the Node in DB
func (n *Node) Save() {

}

// FindNode searches a Node by ID
func FindNode(id string) (*Node, error) {
	var n Node
	var err error

	// TODO: Add logic

	return &n, err
}

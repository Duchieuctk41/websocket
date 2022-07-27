package common

// Message presents a message package that is using to communicate between peers
type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

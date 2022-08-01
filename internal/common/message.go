package common

// Message presents a message package that is using to communicate between peers
type Message struct {
	EventName    string      `json:"eventName"`
	EventPayload interface{} `json:"eventPayload"`
}

package edge_tts

import "fmt"

// UnknownResponseError is raised when an unknown response is received from the server.
type UnknownResponseError struct {
	Message string
}

func (e *UnknownResponseError) Error() string {
	return fmt.Sprintf("Unknown Response: %s", e.Message)
}

// UnexpectedResponseError is raised when an unexpected response is received from the server.
type UnexpectedResponseError struct {
	Message string
}

func (e *UnexpectedResponseError) Error() string {
	return fmt.Sprintf("Unexpected Response: %s", e.Message)
}

// NoAudioReceivedError is raised when no audio is received from the server.
type NoAudioReceivedError struct {
	Message string
}

func (e *NoAudioReceivedError) Error() string {
	return fmt.Sprintf("No Audio Received: %s", e.Message)
}

// WebSocketError is raised when a WebSocket error occurs.
type WebSocketError struct {
	Message string
}

func (e *WebSocketError) Error() string {
	return fmt.Sprintf("WebSocket Error: %s", e.Message)
}

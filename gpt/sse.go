package gpt

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type ChatCompletionChunkEvent struct {
	ID      string        `json:"id"`
	Object  string        `json:"object"`
	Created int           `json:"created"`
	Model   Model         `json:"model"`
	Choices []EventChoice `json:"choices"`
}

type EventChoice struct {
	Delta        Delta   `json:"delta"`
	Index        int     `json:"index"`
	FinishReason *string `json:"finish_reason"` // this is a pointer because the value can be null
}

type Delta struct {
	Role    Role    `json:"role"`
	Content *string `json:"content"`
}

const (
	dataPrefix = "data: "
	doneEvent  = "[DONE]"
)

func parseEventData(data []byte) []byte {
	data = bytes.TrimLeft(data, dataPrefix)
	data = bytes.TrimRight(data, "\n")
	return data
}

func isDoneEvent(data []byte) bool {
	return string(data) == doneEvent
}

func parseEvent(data []byte) (*ChatCompletionChunkEvent, error) {
	event := ChatCompletionChunkEvent{}
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, fmt.Errorf("unmarshal %#v (%s): %w", data, string(data), err)
	}
	return &event, nil
}

package gpt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrimEventData(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		out  []byte
	}{
		{
			name: "ok",
			data: []byte(`data: {"choices": [{"index": 3}]}\n`),
			out:  []byte(`{"choices": [{"index": 3}]}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := parseEventData(tt.data)
			require.Equal(t, tt.out, out)
		})
	}
}

func TestIsDoneEvent(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		out  bool
	}{
		{
			name: "not ok",
			data: nil,
			out:  false,
		},
		{
			name: "ok",
			data: []byte(doneEvent),
			out:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := isDoneEvent(tt.data)
			require.Equal(t, tt.out, out)
		})
	}
}

func TestParseEvent(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		out  *ChatCompletionChunkEvent
		err  error
	}{
		{
			name: "ok",
			data: []byte(`{"choices": [{"index": 3}]}`),
			out: &ChatCompletionChunkEvent{
				Choices: []EventChoice{
					{
						Index: 3,
					},
				},
			},
		},
		{
			name: "ok - actual event",
			data: []byte(`{"id":"chatcmpl-7FkMLWcP8lxIQB5zwefYsXVkAMmoO","object":"chat.completion.chunk","created":1683987481,"model":"gpt-3.5-turbo-0301","choices":[{"delta":{"role":"assistant"},"index":0,"finish_reason":null}]}`),
			out: &ChatCompletionChunkEvent{
				ID:      "chatcmpl-7FkMLWcP8lxIQB5zwefYsXVkAMmoO",
				Object:  "chat.completion.chunk",
				Created: 1683987481,
				Model:   GPT35Turbo0301,
				Choices: []EventChoice{
					{
						Delta: Delta{
							Role: Assistant,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := parseEvent(tt.data)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.out, out)
		})
	}
}

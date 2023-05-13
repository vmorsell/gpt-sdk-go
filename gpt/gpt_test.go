package gpt

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func newClient(t *testing.T) Client {
	_ = godotenv.Load("../.env")

	apiKey := os.Getenv("GPT_API_KEY")
	if apiKey == "" {
		t.Fatalf("api key missing")
	}

	config := NewConfig().WithAPIKey(apiKey)
	return NewClient(config)
}

func TestChatCompletionStream(t *testing.T) {
	tests := []struct {
		in     ChatCompletionInput
		tokens []string
		err    error
	}{
		{
			in: ChatCompletionInput{
				Messages: []Message{
					{
						Role:    RoleUser,
						Content: `Please reply with exactly the text "Hello, World.". Nothing more, nothing less.`,
					},
				},
				Stream: true,
			},
			tokens: []string{"Hello", ",", " World", "."},
		},
	}

	for _, tt := range tests {
		client := newClient(t)
		ch := make(chan *ChatCompletionChunkEvent)

		go func() {
			err := client.ChatCompletionStream(tt.in, ch)
			require.Equal(t, tt.err, err)
		}()

		tokens := []string{}
		for {
			ev, ok := <-ch
			if !ok {
				break
			}

			if ev.Choices[0].Delta.Content == nil {
				continue
			}

			tokens = append(tokens, *ev.Choices[0].Delta.Content)
		}

		require.Equal(t, tt.tokens, tokens)
	}
}

package gpt

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func newClient(t *testing.T) *Client {
	_ = godotenv.Load("../.env")

	apiKey := os.Getenv("GPT_API_KEY")
	if apiKey == "" {
		t.Fatalf("api key missing")
	}

	config := NewConfig().WithAPIKey(apiKey)
	return NewClient(config)
}

func TestChatCompletion(t *testing.T) {
	tests := []struct {
		in      ChatCompletionInput
		choices []Choice
		err     error
	}{
		{
			in: ChatCompletionInput{
				Messages: []Message{
					{
						Role:    RoleUser,
						Content: "Please reply with exactly the text \"Hello, World.\". Nothing more, nothing less.",
					},
				},
			},
			choices: []Choice{
				{
					Index: 0,
					Message: Message{
						Role:    RoleAssistant,
						Content: "Hello, World.",
					},
					FinishReason: "stop",
				},
			},
		},
	}

	for _, tt := range tests {
		client := newClient(t)
		res, err := client.ChatCompletion(tt.in)
		require.Equal(t, tt.err, err)
		require.Equal(t, tt.choices, res.Choices)
	}
}

func ExampleClient_ChatCompletion() {
	config := NewConfig().WithAPIKey("xyz")
	client := NewClient(config)

	in := ChatCompletionInput{
		Messages: []Message{
			{
				Role:    RoleSystem,
				Content: "You are a helpful assistant.",
			},
			{
				Role:    RoleUser,
				Content: `Please reply with the text "Hello, World!". Nothing else.`,
			},
		},
	}
	res, err := client.ChatCompletion(in)
	if err != nil {
		panic(err)
	}

	if res.Choices == nil {
		panic("no choices")
	}

	fmt.Println(res.Choices[0].Message.Content)
	// Output: Hello, World!
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

func ExampleClient_ChatCompletionStream() {
	client := NewClient(NewConfig().WithAPIKey("xyz"))

	in := ChatCompletionInput{
		Messages: []Message{
			{
				Role:    RoleUser,
				Content: `Please reply with the text "Hello, World!". Nothing else.`,
			},
		},
		Stream: true,
	}

	ch := make(chan *ChatCompletionChunkEvent)
	go func() {
		if err := client.ChatCompletionStream(in, ch); err != nil {
			panic(err)
		}
	}()

	for {
		ev, ok := <-ch
		if !ok {
			break
		}

		if ev.Choices[0].Delta.Content == nil {
			continue
		}

		fmt.Print(*ev.Choices[0].Delta.Content)
	}
	// Output: Hello, World!
}

package gpt

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func apiKey() string {
	return os.Getenv("OPENAI_API_KEY")
	}

func TestChatCompletion(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
}

	tests := []struct {
		in      ChatCompletionInput
		choices []Choice
		err     error
	}{
		{
			in: ChatCompletionInput{
				Messages: []Message{
					{
						Role:    User,
						Content: "Please reply with exactly the text \"Hello, World.\". Nothing more, nothing less.",
					},
				},
			},
			choices: []Choice{
				{
					Index: 0,
					Message: Message{
						Role:    Assistant,
						Content: "Hello, World.",
					},
					FinishReason: "stop",
				},
			},
		},
	}

	for _, tt := range tests {
		client := NewClient(NewConfig().WithAPIKey(apiKey()))
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
				Role:    System,
				Content: "You are a helpful assistant.",
			},
			{
				Role:    User,
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
}

func TestChatCompletionStream(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	tests := []struct {
		in     ChatCompletionInput
		tokens []string
		err    error
	}{
		{
			in: ChatCompletionInput{
				Messages: []Message{
					{
						Role:    User,
						Content: `Please reply with exactly the text "Hello, World.". Nothing more, nothing less.`,
					},
				},
				Stream: true,
			},
			tokens: []string{"Hello", ",", " World", "."},
		},
	}

	for _, tt := range tests {
		client := NewClient(NewConfig().WithAPIKey(apiKey()))

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
				Role:    User,
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
}

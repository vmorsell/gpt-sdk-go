package gpt_test

import (
	"fmt"
	"log"

	"github.com/vmorsell/openai-gpt-sdk-go/gpt"
)

func Example_helloWorld() {
	apiKey := ""

	config := gpt.NewConfig().WithAPIKey(apiKey)
	client := gpt.NewClient(config)

	msg := `Can you give me a "Hello, World!"?`
	in := gpt.ChatCompletionInput{
		Messages: []gpt.Message{
			{
				Role:    gpt.System,
				Content: "You are an assistant that speaks like Shakespeare.",
			},
			{
				Role:    gpt.User,
				Content: msg,
			},
		},
	}
	fmt.Printf("User: %s\n", msg)

	res, err := client.ChatCompletion(in)
	if err != nil {
		log.Fatalf("chat completion: %v", err)
	}

	if len(res.Choices) == 0 {
		log.Fatalf("Got 0 choices in the response. This is unexpected.")
	}

	fmt.Printf("ChatGPT: %s\n", res.Choices[0].Message.Content)
}

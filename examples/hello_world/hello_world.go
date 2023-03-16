package main

import (
	"fmt"
	"log"

	"github.com/vmorsell/gpt-sdk-go/gpt"
)

const (
	apiKey = ""
)

func main() {
	config := gpt.NewConfig().WithAPIKey(apiKey)
	client := gpt.NewClient(config)

	msg := "Hi, ChatGPT! Can you give me a \"Hello, World!\"?"
	in := gpt.ChatCompletionInput{
		Messages: []gpt.Message{
			{
				Role:    gpt.RoleUser,
				Content: msg,
			},
		},
	}
	fmt.Printf("User: %s\n", msg)

	res, err := client.ChatCompletion(in)
	if err != nil {
		log.Fatalf("chat completion: %v", err)
	}

	fmt.Printf("ChatGPT: %s\n", res.Choices[0].Message.Content)
}

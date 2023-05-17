package gpt_test

import (
	"fmt"

	"github.com/vmorsell/openai-gpt-sdk-go/gpt"
)

func Example_helloWorld() {
	apiKey := ""
	assistantName := "ShakespeareGPT"

	config := gpt.NewConfig().WithAPIKey(apiKey)
	client := gpt.NewClient(config)

	msg := `Can you give me a "Hello, World!"?`
	in := gpt.ChatCompletionInput{
		Messages: []gpt.Message{
			{
				Role:    gpt.System,
				Content: fmt.Sprintf("You are %s, a helpful assistant that speaks like Shakespeare.", assistantName),
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
		panic(err)
	}

	if len(res.Choices) == 0 {
		panic("Got 0 choices in the response. This is unexpected.")
	}

	fmt.Printf("%s: %s\n", assistantName, res.Choices[0].Message.Content)
}

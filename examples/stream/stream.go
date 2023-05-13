package main

import (
	"fmt"

	"github.com/vmorsell/openai-gpt-sdk-go/gpt"
)

var (
	apiKey = ""
)

func main() {
	config := gpt.NewConfig().WithAPIKey(apiKey)
	client := gpt.NewClient(config)

	msg := `Can you write a haiku about the phrase "Hello, World!"?`
	fmt.Printf("User: %s\n", msg)

	ch := make(chan *gpt.ChatCompletionChunkEvent)
	go func() {
		in := gpt.ChatCompletionInput{
			Messages: []gpt.Message{
				{
					Role:    gpt.RoleSystem,
					Content: "You are an helpful assistant.",
				},
				{
					Role:    gpt.RoleUser,
					Content: msg,
				},
			},
			Stream: true,
		}
		err := client.ChatCompletionStream(in, ch)
		if err != nil {
			panic(fmt.Sprintf("chat completion stream: %v  ", err))
		}
	}()

	fmt.Printf("ChatGPT: ")
	for {
		ev, ok := <-ch
		if !ok {
			break
		}

		if len(ev.Choices) == 0 {
			continue
		}

		if ev.Choices[0].Delta.Content != nil {
			fmt.Print(*ev.Choices[0].Delta.Content)
		}
	}
	fmt.Println()
}

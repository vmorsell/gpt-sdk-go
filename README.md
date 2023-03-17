# OpenAI GPT SDK for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/vmorsell/gpt-sdk-go.svg)](https://pkg.go.dev/github.com/vmorsell/gpt-sdk-go/gpt)

## Getting started

### Installing

    go get github.com/vmorsell/gpt-sdk-go

### Usage

```go
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

    msg := "Can you give me a \"Hello, World!\"?"
    in := gpt.ChatCompletionInput{
        Messages: []gpt.Message{
            {
				Role:    gpt.RoleSystem,
				Content: "You are an assistant that speaks like Shakespeare.",
			},
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

	if len(res.Choices) == 0 {
		log.Fatalf("Got 0 choices in the response. This is unexpected.")
	}

    fmt.Printf("ChatGPT: %s\n", res.Choices[0].Message.Content)
}
```

### How can I get an API key?

Visit https://platform.openai.com

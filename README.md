# OpenAI GPT SDK for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/vmorsell/openai-gpt-sdk-go.svg)](https://pkg.go.dev/github.com/vmorsell/openai-gpt-sdk-go/gpt)

## Getting started

### Installing

    go get github.com/vmorsell/openai-gpt-sdk-go

### Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/vmorsell/openai-gpt-sdk-go/gpt"
)

const (
    apiKey = ""
)

func main() {
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
        panic(err)
    }

    if len(res.Choices) == 0 {
        panic("Got 0 choices in response. This is unexpected")
    }

    fmt.Printf("ShakespeareGPT: %s\n", res.Choices[0].Message.Content)
}
```

### How can I get an API key?

Visit https://platform.openai.com

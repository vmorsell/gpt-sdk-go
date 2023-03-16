/*
Package gpt provides a client for using the OpenAI Chat API.

# Usage

	import "github.com/vmorsell/gpt-sdk-go/gpt"

Initialize a new client.

	config := gpt.NewConfig.WithAPIKey("xyz")
	client := gpt.NewClient(config)

You are now ready to make calls to the API.

	res, err := client.ChatCompletion(gpt.ChatCompletionInput{
		Messages: []gpt.Message{
			{
				Role:    gpt.RoleUser,
				Content: "Hello, World!",
			},
		},
	}

)
*/
package gpt

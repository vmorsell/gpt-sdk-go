// Package gpt is a SDK for OpenAI's API.
package gpt

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	// OpenAI API endpoint.
	DefaultEndpoint = "https://api.openai.com/v1"

	// GPT model to use unless specified in the request.
	DefaultModel = GPT35Turbo

	jsonMIME = "application/json"
)

// Config provides configuration to a client instance.
type Config struct {
	// OpenAI API key.
	APIKey string

	// OpenAI API endpoint. You probably don't want to use this.
	Endpoint string
}

// NewConfig returns a new Config pointer that can be chained with builder
// methods to set multiple configuration values inline without using pointers.
func NewConfig() *Config {
	return &Config{
		Endpoint: DefaultEndpoint,
	}
}

// WithAPIKey sets a config APIKey value returning a Config pointer for
// chaining.
func (c *Config) WithAPIKey(apiKey string) *Config {
	c.APIKey = apiKey
	return c
}

// WithEndpoint sets a config Endpoint value returning a Config pointer for
// chaining.
func (c *Config) WithEndpoint(endpoint string) *Config {
	c.Endpoint = endpoint
	return c
}

type Client struct {
	Config *Config
}

// New creates a new config instance of the OpenAI client. If additional
// configuration is needed for the client instance use the optional aws.Config
// parameter to add your extra config.
func NewClient(config *Config) *Client {
	return &Client{
		Config: config,
	}
}

// makeCall makes a call to the OpenAI API.
func (c *Client) makeCall(httpPath string, payload interface{}, out interface{}) error {
	if payload == nil {
		return fmt.Errorf("empty payload")
	}

	if out == nil {
		return fmt.Errorf("missing return type")
	}

	url := strings.Join([]string{c.Config.Endpoint, httpPath}, "")

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return fmt.Errorf("encode payload: %w", err)
	}

	httpClient := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	req.Header.Add("Content-Type", jsonMIME)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Config.APIKey))

	res, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("post: %w", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return fmt.Errorf("http error %d: %s", res.StatusCode, body)
	}

	if err := json.Unmarshal(body, &out); err != nil {
		return fmt.Errorf("unmarshal %#v (%s): %w", body, string(body), err)
	}

	return nil
}

// makeCallWithResponseStream makes a call to the OpenAI API with tokens
// returned as server-sent events. The output channel is closed when the
// last token is sent.
func (c *Client) makeCallWithResponseStream(httpPath string, payload interface{}, out chan *ChatCompletionChunkEvent) error {
	if httpPath == "" {
		return fmt.Errorf("missing path")
	}

	if payload == nil {
		return fmt.Errorf("missing payload")
	}

	if out == nil {
		return fmt.Errorf("missing output channel")
	}

	url := strings.Join([]string{c.Config.Endpoint, httpPath}, "")

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return fmt.Errorf("encode payload: %w", err)
	}

	httpClient := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	req.Header.Add("Content-Type", jsonMIME)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Config.APIKey))

	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Accept", "text/event-stream")
	req.Header.Add("Connection", "keep-alive")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("post: %w", err)
	}
	defer resp.Body.Close()

	r := bufio.NewReader(resp.Body)
	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("read bytes: %w", err)
		}

		data := parseEventData(line)

		// Sometimes we might get empty lines.
		if data == nil {
			continue
		}

		// Stop if we've got the [DONE] event.
		if isDoneEvent(data) {
			break
		}

		event, err := parseEvent(data)
		if err != nil {
			return fmt.Errorf("parse event: %w", err)
		}

		if len(event.Choices) == 0 {
			return fmt.Errorf("got no choices - this is unexpected")
		}

		out <- event
	}

	close(out)
	return nil
}

const (
	// Path to Chat Completions API endpoint.
	chatCompletionsPath = "/chat/completions"
)

// ChatCompletion implements the chat completion API method.
func (c *Client) ChatCompletion(in ChatCompletionInput) (*ChatCompletionOutput, error) {
	if in.Stream {
		return nil, fmt.Errorf("use ChatCompletionStream method instead to stream return data")
	}

	if in.Model == "" {
		in.Model = DefaultModel
	}

	out := ChatCompletionOutput{}

	if err := c.makeCall(chatCompletionsPath, in, &out); err != nil {
		return nil, fmt.Errorf("make call: %w", err)
	}

	return &out, nil
}

// ChatCompletionStream implements the chat completion API method with a
// response stream. The output channel is closed when the last event has been
// sent.
func (c *Client) ChatCompletionStream(in ChatCompletionInput, out chan *ChatCompletionChunkEvent) error {
	if !in.Stream {
		return fmt.Errorf("use ChatCompletion method instead")
	}

	if in.Model == "" {
		in.Model = DefaultModel
	}

	if err := c.makeCallWithResponseStream(chatCompletionsPath, in, out); err != nil {
		return fmt.Errorf("make call with stream response: %w", err)
	}

	return nil
}

// GPT model names.
type Model string

const (
	// GPT-4.
	// Note: You currently need special grant to use these (2023-03-16).
	GPT4        Model = "gpt-4"
	GPT40314    Model = "gpt-4-0314"
	GPT432k     Model = "gpt-4-32k"
	GPT432k0314 Model = "gpt-4-32k-0314"

	// GPT-3.5.
	GPT35Turbo     Model = "gpt-3.5-turbo"
	GPT35Turbo0301 Model = "gpt-3.5-turbo-0301"
)

// Message roles.
type Role string

const (
	// The system message helps set the behavior of the assistant.
	// The system message can for example be "You are a helpful
	// assistant."
	System Role = "system"

	// The assistant messages help store prior responses. They
	// can also be written by a developer to help give examples
	// of desired behavior.
	Assistant Role = "assistant"

	// The user messages help instruct the assistant. They can
	// be generated by the end users of an application, or set
	// by a developer as an instruction.
	User Role = "user"
)

// ChatCompletionInput is the input to a ChatCompletion call.
type ChatCompletionInput struct {
	Model    Model     `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

// ChatCompletionOutput is the output to a ChatCompletion call.
type ChatCompletionOutput struct {
	Id      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice is a single returned suggestion of what the next chat
// message could be.
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Message represents a message.
type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

// Usage holds token usage reporting from ChatGPT.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

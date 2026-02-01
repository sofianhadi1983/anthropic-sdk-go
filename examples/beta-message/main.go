package main

import (
	"context"
	"fmt"

	"github.com/sofianhadi1983/anthropic-sdk-go"
)

func main() {
	client := anthropic.NewClient()

	content := "What is the capital of France? Answer in one word."

	println("[user]: " + content)

	// Using the Beta Messages API with beta features enabled via the Betas param
	message, err := client.Beta.Messages.New(context.TODO(), anthropic.BetaMessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.BetaMessageParam{
			anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock(content)),
		},
		Model: anthropic.ModelClaudeSonnet4_5_20250929,
		// Optional: Enable specific beta features
		// Betas: []anthropic.AnthropicBeta{
		// 	anthropic.AnthropicBetaOutput128k2025_02_19,
		// },
	})
	if err != nil {
		panic(err)
	}

	// Access response content
	for _, block := range message.Content {
		switch variant := block.AsAny().(type) {
		case anthropic.BetaTextBlock:
			fmt.Printf("[assistant]: %s\n", variant.Text)
		}
	}

	// Print usage info
	fmt.Printf("\nUsage: input=%d, output=%d tokens\n",
		message.Usage.InputTokens, message.Usage.OutputTokens)
}

package main

import (
	"context"
	"fmt"

	"github.com/sofianhadi1983/anthropic-sdk-go"
)

func main() {
	client := anthropic.NewClient()

	content := "Write a haiku about programming."

	println("[user]: " + content)
	print("[assistant]: ")

	// Using the Beta Messages API with streaming
	stream := client.Beta.Messages.NewStreaming(context.TODO(), anthropic.BetaMessageNewParams{
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

	// Accumulate the full message while streaming
	message := anthropic.BetaMessage{}
	for stream.Next() {
		event := stream.Current()

		// Accumulate events into the final message
		if err := message.Accumulate(event); err != nil {
			fmt.Printf("error accumulating event: %v\n", err)
			continue
		}

		// Handle streaming events
		switch eventVariant := event.AsAny().(type) {
		case anthropic.BetaRawContentBlockDeltaEvent:
			switch deltaVariant := eventVariant.Delta.AsAny().(type) {
			case anthropic.BetaTextDelta:
				print(deltaVariant.Text)
			}
		case anthropic.BetaRawMessageDeltaEvent:
			// Handle stop sequence if present
			if eventVariant.Delta.StopSequence != "" {
				print(eventVariant.Delta.StopSequence)
			}
		}
	}

	println()

	if stream.Err() != nil {
		panic(stream.Err())
	}

	// Print usage info from accumulated message
	fmt.Printf("\nUsage: input=%d, output=%d tokens\n",
		message.Usage.InputTokens, message.Usage.OutputTokens)
}

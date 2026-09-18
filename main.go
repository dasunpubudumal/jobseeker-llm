package main

// TODO - Use: https://platform.claude.com/docs/en/cli-sdks-libraries/sdks/go
// Use https://developer.adzuna.com/admin/access_details for the API
//
// For reference, use https://platform.claude.com/docs/en/agents-and-tools/tool-use/overview

import (
	"context"
	"fmt"
	"log"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/joho/godotenv"
	"github.com/dasunpubudumal/jobseeker-llm/adzuna_client"
)

func load_env() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file.")
	}
}

func main() {

	load_env()

	adzuna := adzuna_client.AdzunaClient{
		APIKey: "",
		BaseUrl: "",
	}

	adzuna.GetJobsForAJobType("JavaScript Developer")

	client := anthropic.NewClient()

	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(
				anthropic.NewTextBlock("What is a black hole?"),
				),
		},
		Model: anthropic.ModelClaudeSonnet4_5,
	})

	if err != nil {
		panic(err.Error())
	}

	for _, block := range message.Content {
		if textBlock, ok := block.AsAny().(anthropic.TextBlock); ok {
			fmt.Println(textBlock.Text)
		}
	}
}

package main

// Used: https://platform.claude.com/docs/en/cli-sdks-libraries/sdks/go
// Used https://developer.adzuna.com/admin/access_details for the API
//
// For reference, use https://platform.claude.com/docs/en/agents-and-tools/tool-use/overview

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	adzuna "github.com/dasunpubudumal/jobseeker-llm/adzuna_client"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

type ToolUseType struct {
	JobType string `json:"job_type"`
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file.")
	}
}

func run(prompt string) {
	restyClient := resty.New()
	adzuna := adzuna.AdzunaClient{
		AppID:      os.Getenv("ADZUNA_CLIENT_ID"),
		AppKey:     os.Getenv("ADZUNA_CLIENT_SECRET"),
		BaseURL:    "https://api.adzuna.com/v1/api",
		HTTPClient: restyClient,
	}

	client := anthropic.NewClient()
	context := context.Background()

	tools := []anthropic.ToolUnionParam{
		{OfTool: &anthropic.ToolParam{
			Name:        "get_jobs_for_a_type",
			Description: anthropic.String("Get jobs from a given job type."),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]any{
					"job_type": map[string]any{
						"type":        "string",
						"description": "Type of the job, e.g., JavaScript Developer",
					},
				},
				Required: []string{"job_type"},
			},
		}},
	}

	toolChoice := anthropic.ToolChoiceUnionParam{
		OfAuto: &anthropic.ToolChoiceAutoParam{DisableParallelToolUse: anthropic.Bool(true)},
	}
	messages := []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
	}

	response, err := client.Messages.New(context, anthropic.MessageNewParams{
		Model:      anthropic.ModelClaudeSonnet4_5,
		MaxTokens:  1024,
		Tools:      tools,
		ToolChoice: toolChoice,
		Messages:   messages,
	})
	if err != nil {
		log.Fatal(err)
	}

	var toolUse anthropic.ContentBlockUnion
	for _, block := range response.Content {
		if block.Type == "tool_use" {
			toolUse = block
			break
		}
	}

	fmt.Printf("Claude called %s with %s\n", toolUse.Name, string(toolUse.Input))

	// Now, run the tool.
	args, err := toolUse.Input.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	toolUseStruct := &ToolUseType{}
	json.Unmarshal(args, toolUseStruct)

	jobs := adzuna.GetJobsForAJobType(toolUseStruct.JobType)
	var assistantContent []anthropic.ContentBlockParamUnion
	for _, block := range response.Content {
		assistantContent = append(assistantContent, block.ToParam())
	}
	messages = append(
		messages,
		anthropic.NewAssistantMessage(assistantContent...),
		anthropic.NewUserMessage(anthropic.NewToolResultBlock(
			toolUse.ID,
			jobs.AsJSONString(),
			false,
		)),
	)
	followup, err := client.Messages.New(context, anthropic.MessageNewParams{
		Model:      anthropic.ModelClaudeSonnet4_5,
		MaxTokens:  1024,
		Tools:      tools,
		ToolChoice: toolChoice,
		Messages:   messages,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Claude uses the result to answer the original question
	for _, block := range followup.Content {
		if block.Type == "text" {
			fmt.Println("The LLM responded with: ")
			fmt.Println(block.Text)
		}
	}
}

func main() {
	loadEnv()

	for {
		fmt.Println("Input text: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		err := scanner.Err()
		if err != nil {
			log.Fatal(err)
		}

		text := scanner.Text()

		if len(text) == 0 {
			fmt.Println("Exiting the program.")
			break
		}

		run(text)
	}
}

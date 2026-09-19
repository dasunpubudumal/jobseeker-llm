package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/dasunpubudumal/jobseeker-llm/adzuna"
)

type ToolUseType struct {
	JobType string `json:"job_type"`
}

type CaludeClient struct{}

func (c *CaludeClient) Invoke(prompt string, adzunaClient adzuna.AdzunaClient) LLMResponse {
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

	log.Println("Claude called %s with %s\n", toolUse.Name, string(toolUse.Input))

	// Now, run the tool.
	args, err := toolUse.Input.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	toolUseStruct := &ToolUseType{}
	json.Unmarshal(args, toolUseStruct)

	jobs := adzunaClient.GetJobsForAJobType(toolUseStruct.JobType)
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

	var resp string

	// Claude uses the result to answer the original question
	for _, block := range followup.Content {
		if block.Type == "text" {
			fmt.Println("The LLM responded with: ")
			resp = block.Text
		}
	}

	return LLMResponse{
		response: resp,
	}
}

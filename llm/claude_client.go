package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/dasunpubudumal/jobseeker-llm/adzuna"
)

type ToolUseType struct {
	JobType  string `json:"job_type"`
	Location string `json:"location"`
}

type CaludeClient struct{}

func (c *CaludeClient) Invoke(prompt string, adzunaClient adzuna.AdzunaClient) (LLMResponse, error) {
	client := anthropic.NewClient()
	context := context.Background()

	tools := []anthropic.ToolUnionParam{
		{
			OfTool: &anthropic.ToolParam{
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
			},
		},
		{
			OfTool: &anthropic.ToolParam{
				Name:        "get_jobs_for_a_type_and_location",
				Description: anthropic.String("Get jobs from a given job type and based on a given location."),
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: map[string]any{
						"job_type": map[string]any{
							"type":        "string",
							"description": "Type of the job, e.g., JavaScript Developer",
						},
						"location": map[string]any{
							"type":        "string",
							"description": "Location of the job e.g., San Fransisco",
						},
					},
					Required: []string{"location"},
				},
			},
		},
	}

	toolChoice := anthropic.ToolChoiceUnionParam{
		OfAuto: &anthropic.ToolChoiceAutoParam{DisableParallelToolUse: anthropic.Bool(true)},
	}

	PROMPT := fmt.Sprintf(
		`	You are a Recruitment Specialist.

		Your task is to find job opportunities for the user, based on a query
		given by the user.

		User Query: %s.

		Please provide the response according to the following template:

		- Title of the job
		- Minimum Salary
		- Job Description
		- Company Name
		- Contract Time
		- URL
	`,
		prompt,
	)

	messages := []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock(PROMPT)),
	}

	response, err := client.Messages.New(context, anthropic.MessageNewParams{
		Model:      anthropic.ModelClaudeSonnet4_5,
		MaxTokens:  1024,
		Tools:      tools,
		ToolChoice: toolChoice,
		Messages:   messages,
	})
	if err != nil {
		return LLMResponse{}, err
	}

	var toolUse anthropic.ContentBlockUnion
	for _, block := range response.Content {
		if block.Type == "tool_use" {
			toolUse = block
			break
		}
	}

	log.Printf("Claude called %s with %s\n", toolUse.Name, string(toolUse.Input))

	// Now, run the tool.
	args, err := toolUse.Input.MarshalJSON()
	if err != nil {
		return LLMResponse{}, err
	}
	toolUseStruct := &ToolUseType{}
	err = json.Unmarshal(args, toolUseStruct)
	if err != nil {
		return LLMResponse{}, err
	}

	var jobs adzuna.AdzunaResponse
	var adzerr error

	switch toolUse.Name {
	case "get_jobs_for_a_type_and_location":
		jobs, adzerr = adzunaClient.GetJobsForAJobTypeAndLocation(toolUseStruct.JobType, toolUseStruct.Location)
	case "get_jobs_for_a_type":
		jobs, adzerr = adzunaClient.GetJobsForAJobType(toolUseStruct.JobType)
	default:
		return LLMResponse{}, errors.New("internal error from one of the downstream tool calls")
	}

	if adzerr != nil {
		return LLMResponse{}, err
	}

	var assistantContent []anthropic.ContentBlockParamUnion
	for _, block := range response.Content {
		assistantContent = append(assistantContent, block.ToParam())
	}
	jobsJSON, err := jobs.AsJSONString()
	if err != nil {
		return LLMResponse{}, err
	}
	messages = append(
		messages,
		anthropic.NewAssistantMessage(assistantContent...),
		anthropic.NewUserMessage(anthropic.NewToolResultBlock(
			toolUse.ID,
			jobsJSON,
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
		return LLMResponse{}, err
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
	}, nil
}

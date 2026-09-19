package main

// Used: https://platform.claude.com/docs/en/cli-sdks-libraries/sdks/go
// Used https://developer.adzuna.com/admin/access_details for the API
//
// For reference, use https://platform.claude.com/docs/en/agents-and-tools/tool-use/overview

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dasunpubudumal/jobseeker-llm/adzuna"
	"github.com/dasunpubudumal/jobseeker-llm/llm"
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

func main() {
	loadEnv()
	var model string
	var llmClient llm.LLMClient

	model = os.Getenv("MODEL")

	if strings.Contains(model, "claude") {
		llmClient = &llm.CaludeClient{}
	} else {
		log.Fatal("The system has not been configured to use any LLMs.")
	}

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

		restyClient := resty.New()
		adzuna := adzuna.AdzunaClient{
			AppID:      os.Getenv("ADZUNA_CLIENT_ID"),
			AppKey:     os.Getenv("ADZUNA_CLIENT_SECRET"),
			BaseURL:    "https://api.adzuna.com/v1/api",
			HTTPClient: restyClient,
		}

		resp, err := llmClient.Invoke(text, adzuna)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resp)
	}
}

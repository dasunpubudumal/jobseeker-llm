package adzuna_client

// RESTY Documentation https://github.com/go-resty/resty/blob/v2/README.md

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

type AdzunaClient struct {
	AppID      string
	AppKey     string
	BaseURL    string
	HTTPClient *resty.Client
}

type AdzunaResponse struct {
	Count   int               `json:"count"`
	Results []AdjunaJobResult `json:"results"`
}

type AdjunaJobResult struct {
	SalaryMin         float32                 `json:"salary_min"`
	Title             string                  `json:"title"`
	SalaryIsPredicted string                  `json:"salary_is_predicted"`
	RedirectURL       string                  `json:"redirect_url"`
	Description       string                  `json:"description"`
	Company           AdjunaJobResultCompany  `json:"company"`
	ContractTime      string                  `json:"contract_time"`
	Category          AdjunaJobResultCategory `json:"category"`
	Created           string                  `json:"created"`
}

type AdjunaJobResultCompany struct {
	DisplayName string `json:"display_name"`
}

type AdjunaJobResultCategory struct {
	Label string `json:"label"`
	Tag   string `json:"tag"`
}

func (c *AdzunaClient) GetJobsForAJobType(jobType string) AdzunaResponse {
	var response AdzunaResponse
	_, err := c.HTTPClient.
		R().
		SetResult(&response).
		SetQueryParam("app_id", c.AppID).
		SetQueryParam("app_key", c.AppKey).
		SetQueryParam("what", jobType).
		EnableTrace().
		ExpectContentType("application/json").
		Get(fmt.Sprintf("%s/jobs/gb/search/1", c.BaseURL))
	if err != nil {
		log.Fatal(fmt.Errorf("error in sending the request to Adjuna: %v", err))
	}

	return response
}

/* Package adzuna contains the Adzuna client that connects with Adzuna API */
package adzuna

// RESTY Documentation https://github.com/go-resty/resty/blob/v2/README.md

import (
	"encoding/json"
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

func (r *AdzunaResponse) AsJSONString() (string, error) {
	bytes, err := json.Marshal(r)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
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

func (c *AdzunaClient) GetJobsForAJobTypeAndLocation(jobType string, location string) (AdzunaResponse, error) {
	count := 1
	resultCount := 0
	results := []AdjunaJobResult{}

	for {
		var response AdzunaResponse
		log.Printf("Adzuna was called for the /search/%d page.", count)
		_, err := c.HTTPClient.
			R().
			SetResult(&response).
			SetQueryParam("app_id", c.AppID).
			SetQueryParam("app_key", c.AppKey).
			SetQueryParam("what", jobType).
			SetQueryParam("where", location).
			EnableTrace().
			ExpectContentType("application/json").
			Get(fmt.Sprintf("%s/jobs/gb/search/%d", c.BaseURL, count))
		if len(response.Results) == 0 {
			break
		}
		results = append(
			results, response.Results...,
		)
		// This needs to happen only once.
		resultCount = response.Count
		if err != nil {
			return AdzunaResponse{}, err
		}
		count++
	}

	return AdzunaResponse{Count: resultCount, Results: results}, nil
}

func (c *AdzunaClient) GetJobsForAJobType(jobType string) (AdzunaResponse, error) {
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
		return AdzunaResponse{}, err
	}

	log.Printf("Adzuna was called for %s; responded with %d results", jobType, response.Count)

	return response, nil
}

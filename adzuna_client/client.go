package adzuna_client

// RESTY Documentation https://github.com/go-resty/resty/blob/v2/README.md

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)


type AdzunaClient struct {
	AppId	string
	AppKey string
	BaseUrl	string
	HttpClient *resty.Client
}

func (c *AdzunaClient) GetJobsForAJobType(job_type string) string {
	resp, err := c.HttpClient.
		R().
		SetQueryParam("app_id", c.AppId).
		SetQueryParam("app_key", c.AppKey).
		SetQueryParam("what", job_type).
		EnableTrace().
		Get(fmt.Sprintf("%s/jobs/gb/search/1", c.BaseUrl))

	if err != nil {
		log.Fatal(fmt.Errorf("Error in sending the request to Adjuna: %v", err))
	}

	return resp.String()
}



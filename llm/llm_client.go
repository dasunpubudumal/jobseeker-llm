/* Package llm is contains the interface for downstream LLM clients */
package llm

import "github.com/dasunpubudumal/jobseeker-llm/adzuna"

type LLMResponse struct {
	response string
}

type LLMClient interface {
	Invoke(prompt string, adzunaClient adzuna.AdzunaClient) (LLMResponse, error)
}

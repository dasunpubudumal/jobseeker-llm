# jobseeker-llm

A small Go CLI that lets an LLM search for jobs. You type a request in plain English, Claude decides which tool to call (job type only, or job type plus location), and the program fetches matching listings from the [Adzuna](https://developer.adzuna.com/) API and hands them back to Claude to summarise.

## How it works

1. You enter a prompt (e.g. "Find me JavaScript developer jobs in London").
2. The prompt is wrapped in a "Recruitment Specialist" system-style prompt and sent to Claude along with two tools:
   - `get_jobs_for_a_type` takes a `job_type`.
   - `get_jobs_for_a_type_and_location` takes a `job_type` and a `location`.
3. Claude responds with a tool call, and the program runs the matching Adzuna query:
   - Job type only: fetches the first page of results (`/jobs/gb/search/1`).
   - Job type and location: pages through `/jobs/gb/search/{n}` until a page comes back empty, collecting every result.
4. The results are serialised to JSON and sent back to Claude as a tool result.
5. Claude's final answer is printed to the terminal. It reports how many results were found and lists up to 10, ordered by salary (descending), with title, minimum salary, description, company, contract time and URL.

The loop repeats until you submit an empty line.

## Project layout

```
.
├── main.go                # CLI loop, env loading and LLM client selection
├── adzuna/
│   └── client.go          # Adzuna API client and response types
└── llm/
    ├── llm_client.go      # LLMClient interface and LLMResponse type
    ├── claude_client.go   # Claude implementation (tool definitions and tool-use flow)
    └── ollama_client.go   # Placeholder for a future Ollama client (empty)
```

## Prerequisites

- Go (see `go.mod` for the version)
- An [Anthropic API key](https://platform.claude.com/)
- An [Adzuna](https://developer.adzuna.com/) app ID and app key

## Setup

1. Copy the example environment file and fill in your credentials:

   ```sh
   cp .env_example .env
   ```

   | Variable               | Description                                                                 |
   | ---------------------- | --------------------------------------------------------------------------- |
   | `ANTHROPIC_API_KEY`    | Your Anthropic API key                                                      |
   | `ADZUNA_CLIENT_ID`     | Your Adzuna application ID                                                  |
   | `ADZUNA_CLIENT_SECRET` | Your Adzuna application key                                                 |
   | `MODEL`                | Selects the LLM backend. Must contain `claude` (e.g. `MODEL=claude`)        |

   `MODEL` is not in `.env_example`, so add it to your `.env` yourself. The program exits with "The system has not been configured to use any LLMs." if it is unset or doesn't contain `claude`.

   `.env` is git-ignored, so your keys won't be committed.

2. Install dependencies:

   ```sh
   go mod download
   ```

## Usage

```sh
go run .
```

At the `Input text:` prompt, ask for jobs, for example:

```
Input text:
Are there any JavaScript developer jobs?
```

or, with a location:

```
Input text:
Find me data engineer jobs in Manchester
```

Submit an empty line to exit.

## Notes

- Searches currently target Adzuna's Great Britain (`gb`) endpoint.
- Location searches fetch every page of results, so broad queries can make many Adzuna requests.
- The Claude client uses `claude-sonnet-4-5` (`anthropic.ModelClaudeSonnet4_5`) regardless of the `MODEL` value; `MODEL` only picks the backend.
- Only Claude is supported today. `llm.LLMClient` is the interface a new backend would implement.

## Dependencies

- [anthropic-sdk-go](https://github.com/anthropics/anthropic-sdk-go) for the Claude API
- [resty](https://github.com/go-resty/resty) for HTTP requests to Adzuna
- [godotenv](https://github.com/joho/godotenv) for loading `.env`

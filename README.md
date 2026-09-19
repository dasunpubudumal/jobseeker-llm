# jobseeker-llm

A small Go CLI that lets Claude search for jobs. You type a request in plain English, Claude decides which job type to look up, and the program fetches matching listings from the [Adzuna](https://developer.adzuna.com/) API and hands them back to Claude to summarise.

## How it works

1. You enter a prompt (e.g. "Find me JavaScript developer jobs").
2. Claude is given a `get_jobs_for_a_type` tool and responds with a tool call containing a `job_type`.
3. The program calls Adzuna's search endpoint (`/jobs/gb/search/1`) for that job type.
4. The results are serialised to JSON and sent back to Claude as a tool result.
5. Claude's final answer is printed to the terminal.

The loop repeats until you submit an empty line.

## Project layout

```
.
├── main.go            # CLI loop and Claude tool-use flow
└── adzuna_client/
    └── client.go      # Adzuna API client and response types
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

   | Variable               | Description                    |
   | ---------------------- | ------------------------------ |
   | `ANTHROPIC_API_KEY`    | Your Anthropic API key         |
   | `ADZUNA_CLIENT_ID`     | Your Adzuna application ID     |
   | `ADZUNA_CLIENT_SECRET` | Your Adzuna application key    |

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

Submit an empty line to exit.

## Notes

- Searches currently target Adzuna's Great Britain (`gb`) endpoint and return the first page of results only.
- The program uses `claude-sonnet-4-5` (`anthropic.ModelClaudeSonnet4_5`).

## Dependencies

- [anthropic-sdk-go](https://github.com/anthropics/anthropic-sdk-go) for the Claude API
- [resty](https://github.com/go-resty/resty) for HTTP requests to Adzuna
- [godotenv](https://github.com/joho/godotenv) for loading `.env`

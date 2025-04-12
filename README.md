# flowmeter-llm

Calculates average water flow rate from images of water meters using LLM (Large Language Model) processing.  Compatible with any OpenAI-compatible API endpoint (e.g., LiteLLM, Ollama).  Allows specifying the LLM model for processing.  The `gemini-flash` model has shown the best price/performance ratio in testing.

**Disclaimer:** This project serves a practical purpose, but is also being used as a learning exercise for the Go programming language.

## Build

```bash
cd flowmeter-llm
go build .
```

## Usage

### Get help

```bash
wmeter --help
```

### Scan and ingest new data:

```bash
OPENAI_API_KEY="<<OPENAI_KEY>>" wmeter ingest /path/to/meter/data --endpoint "<<OPEN_AI_COMPATIBLE_BASE_URL>>"
```

* Replace <<OPENAI_KEY>> with your actual OpenAI API key.
* Replace <<OPEN_AI_COMPATIBLE_BASE_URL>> with the base URL of your OpenAI-compatible API.  For example, http://localhost:8000/v1.
     

### Analyze the ingested data:

```bash
./wmeter analyze -n 10
```

* `-n 10` analyzes the last 10 readings. Adjust the number as needed.
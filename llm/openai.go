package llm

import (
	"context"
	// "errors"
	"encoding/json"
	"fmt"
	"log"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/snexus/wmeter/entities"
	"github.com/snexus/wmeter/fs"
)

func GetOpenAIClient(apiKey, baseUrl string) *openai.Client {
	// Create a new OpenAI client with the provided API key
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = baseUrl
	return openai.NewClientWithConfig(config)
}

func DescribeImage(client *openai.Client, imagePath string, prompt string, modelName string) (*entities.MeterReadingResult, error) {
	log.Printf("Analyzing: Image path: %s\n", imagePath)
	base64Image, err := fs.ReadImageFileToBase64(imagePath)

	if err != nil {
		panic(err)
	}

	// Create the result schem
	var result entities.MeterReadingResult
	schema, err := jsonschema.GenerateSchemaForType(result)
	// fmt.Printf("Schema: %s\n", schema)

	if err != nil {
		log.Fatalf("GenerateSchemaForType error: %v", err)
	}

	// Create the request with image content
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			// Model: openai.GPT4oMini,
			Model: modelName,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleUser,
					MultiContent: []openai.ChatMessagePart{
						{
							Type: openai.ChatMessagePartTypeText,
							Text: prompt,
						},
						{
							Type: openai.ChatMessagePartTypeImageURL,
							ImageURL: &openai.ChatMessageImageURL{
								URL: fmt.Sprintf("data:image/jpeg;base64,%s", base64Image),
							},
						},
					},
				},
			},
			Temperature: 0,
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
				JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
					Name:   "water_meter_reading",
					Schema: schema,
					Strict: false,
				},
			},
		},
	)


	if err != nil {
		fmt.Println("There was a problem with chat completion: %v\n", err)
	}

	// err = schema.Unmarshal(resp.Choices[0].Message.Content, &result)
	fmt.Println("Result content: ", resp.Choices[0].Message.Content)
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result)


	if err != nil {
		log.Fatalf("Unmarshal schema error: %v", err)
	}
	return &result, nil
}

package openai

import (
	"context"
	"github.com/sashabaranov/go-openai"
	schema "github.com/sashabaranov/go-openai/jsonschema"
	"os"
)

func GetClient() *openai.Client {
	apiToken := os.Getenv("TOKEN")
	return openai.NewClient(apiToken)
}

func AIFunctionSchema() openai.FunctionDefinition {
	return openai.FunctionDefinition{
		Name: "quiz_generator",
		Parameters: schema.Definition{
			Type: schema.Object,
			Properties: map[string]schema.Definition{
				"questions": {
					Type:        schema.Array,
					Description: "Array of questions",
					Items: &schema.Definition{
						Type: schema.Object,
						Properties: map[string]schema.Definition{
							"title": {
								Type:        schema.String,
								Description: "Question title",
							},
							"options": {
								Type:        schema.Array,
								Description: "Array of options",
								Items: &schema.Definition{
									Type: schema.Object,
									Properties: map[string]schema.Definition{
										"title": {
											Type:        schema.String,
											Description: "Title of the option",
										},
										"is_correct": {
											Type:        schema.Boolean,
											Description: "flag that indicates whether the option is correct",
										},
									},
									Required: []string{"option", "is_correct"},
								},
							},
						},
					},
				},
			},
			Required: []string{"questions"},
		},
	}
}

func GetModels() []string {
	client := GetClient()
	models, err := client.ListModels(context.Background())
	if err != nil {
		panic(err)
	}

	var modelList []string
	for _, model := range models.Models {
		modelList = append(modelList, model.ID)
	}

	return modelList
}

func CreateChatCompletion(request openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	client := GetClient()
	return client.CreateChatCompletion(context.Background(), request)
}

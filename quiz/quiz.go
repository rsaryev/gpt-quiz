package quiz

import (
	"container/list"
	"encoding/json"
	"github.com/sashabaranov/go-openai"
	internal_openai "gpt-quiz/internal/openai"
)

type Context struct {
	Content   string
	Questions []Question
	Queue     *list.List
}

type Option struct {
	Title     string `json:"title"`
	IsCorrect bool   `json:"is_correct"`
}

type Question struct {
	Title   string   `json:"title"`
	Options []Option `json:"options"`
}

type Data struct {
	Questions []Question `json:"questions"`
}

func createQuestion(ctx *Context) (openai.ChatCompletionResponse, error) {
	prompt := `Create questions and answers based on text.` + ctx.Content + `
Restrictions:
- There should be 7 questions and 4 answer options
- In the answers there should be only 1 correct answer and 3 incorrect ones (no more and no less)
- Questions and answers should be according to best practices for knowledge testing
	`

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
	}

	request := openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo0613,
		Messages: messages,
		FunctionCall: &openai.FunctionCall{
			Name: "quiz_generator",
		},
		Functions: []openai.FunctionDefinition{
			internal_openai.AIFunctionSchema(),
		},
		Temperature: 0.7,
	}

	response, err := internal_openai.CreateChatCompletion(request)
	return response, err
}

func getQuizData(context *Context) (Data, error) {
	excludeQuestions := make([]string, 0)
	for _, question := range context.Questions {
		excludeQuestions = append(excludeQuestions, question.Title)
	}
	response, err := createQuestion(context)
	if err != nil {
		return Data{}, err
	}

	args := response.Choices[0].Message.FunctionCall.Arguments

	var quizData Data
	err = json.Unmarshal([]byte(args), &quizData)
	if err != nil {
		return Data{}, err
	}

	return quizData, nil
}

func CreateQuestions(context *Context) error {
	if context.Queue == nil {
		context.Queue = list.New()
	}
	if context.Queue.Len() > 3 {
		return nil
	}

	data, _ := getQuizData(context)
	for _, question := range data.Questions {
		context.Queue.PushBack(question)
	}

	context.Questions = append(context.Questions, data.Questions...)
	return nil
}

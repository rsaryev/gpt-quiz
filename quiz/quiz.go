package quiz

import (
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	internal_openai "gpt-quiz/internal/openai"
	"sync"
)

type Context struct {
	Content string
	Mux     *sync.Mutex
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
	prompt := fmt.Sprintf(`Create questions and answers based on text: %s
Restrictions:
- There should be 10 questions and 4 answer options
- There should be only 1 correct answer and 3 incorrect ones (no more and no less)
- Questions and answers should follow best practices for knowledge testing`, ctx.Content)

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

/**
 * Filter questions by the following criteria (see createQuestion function prompt)
 */
func filterQuestions(questions []Question) []Question {
	var filtered []Question
	for _, question := range questions {
		if len(question.Options) != 4 {
			continue
		}

		correctAnswers := 0
		for _, option := range question.Options {
			if option.IsCorrect {
				correctAnswers++
			}
		}

		if correctAnswers != 1 {
			continue
		}

		filtered = append(filtered, question)
	}

	return filtered
}
func GetQuestions(context *Context) (Data, error) {
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

	quizData.Questions = filterQuestions(quizData.Questions)

	return quizData, nil
}

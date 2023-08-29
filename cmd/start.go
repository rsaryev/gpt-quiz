package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"gpt-quiz/quiz"
	"time"
)

func Start() *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "start",
		Aliases: []string{"s"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("requires text")
			}
			return nil
		},

		Short: "Start quiz",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := start(quiz.Context{
				Content: args[0],
			})
			if err != nil {
				println(err.Error())
			}
			return nil
		},
		Example: "gpt-quiz start \"golang interview questions\"",
	}

	return cmd
}

func getCorrectAnswer(question quiz.Question) string {
	for _, questionOption := range question.Options {
		if questionOption.IsCorrect {
			return questionOption.Title
		}
	}
	return "No found correct answer"
}

func runQuestions(question quiz.Question) {
	options := make([]string, len(question.Options))

	for i, option := range question.Options {
		options[i] = option.Title
	}

	prompt := promptui.Select{
		Label: fmt.Sprintf(question.Title),
		Items: options,
		Templates: &promptui.SelectTemplates{
			Active:   `> {{ . | cyan }}`,
			Inactive: `  {{ . | cyan }}`,
			Selected: `{{ "` + question.Title + `" }}: {{ . | green }}`,
		},
		HideHelp: true,
	}

	_, option, err := prompt.Run()
	if err != nil {
		panic(err)
	}

	answer := getCorrectAnswer(question)
	if option != answer {
		fmt.Println(color.RedString(">") + " " + color.GreenString(answer))
	}

	fmt.Println("")
}

func loopLoadQuestions(context quiz.Context, ch chan quiz.Question) {
	for {
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)

		if len(ch) == 0 {
			s.Suffix = " Loading questions"
			s.Start()
		}

		q, err := quiz.GetQuestions(&context)
		if err != nil {
			panic(err)
		}

		for _, question := range q.Questions {
			if question.Title == "" {
				continue
			}
			ch <- question
		}
		s.Stop()
	}
}

func start(context quiz.Context) error {
	ch := make(chan quiz.Question, 30)

	go loopLoadQuestions(context, ch)
	for question := range ch {
		runQuestions(question)
	}

	return nil
}

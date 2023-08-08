package cmd

import (
	"fmt"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"gpt-quiz/quiz"
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

func start(context quiz.Context) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	isRunning := false

	for {
		if !isRunning {
			isRunning = true
			wg.Add(1)
			go func() {
				mu.Lock()
				quiz.CreateQuestions(&context)
				isRunning = false
				wg.Done()
				mu.Unlock()
			}()
		}

		if context.Queue == nil || context.Queue.Len() == 0 {
			s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
			s.Suffix = " Loading..."
			s.Start()
			wg.Wait()
			s.Stop()
		} else {
			question := context.Queue.Front().Value.(quiz.Question)
			context.Queue.Remove(context.Queue.Front())
			var options []string
			for _, option := range question.Options {
				options = append(options, option.Title)
			}

			prompt := promptui.Select{
				Label: fmt.Sprintf(question.Title),
				Items: append(options, "Skip"),
				Templates: &promptui.SelectTemplates{
					Active:   `> {{ . | cyan }}`,
					Inactive: `  {{ . | cyan }}`,
					Selected: `{{ "` + question.Title + `" }}: {{ . | green }}`,
				},
				HideHelp: true,
			}

			_, option, err := prompt.Run()
			if err != nil {
				return err
			}

			if option == "Skip" {
				fmt.Println(color.YellowString(">") + " " + color.GreenString(getCorrectAnswer(question)))
				continue
			}

			answer := getCorrectAnswer(question)
			if option != answer {
				fmt.Println(color.RedString(">") + " " + color.GreenString(answer))
			}

			fmt.Println("")
		}
	}
}

package config

import (
	"encoding/json"
	"github.com/manifoldco/promptui"
	"gpt-quiz/internal/openai"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	TOKEN string `json:"token"`
	MODEL string `json:"model"`
}

func (c *Config) write(filename string) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	c.setEnvVariables()
	return nil
}

func (c *Config) configureToken(filename string) {
	if c.TOKEN != "" {
		return
	}

	prompt := promptui.Prompt{
		Label: "🔑 Please enter your token",
	}

	result, err := prompt.Run()
	if err != nil {
		panic(err)
	}

	c.TOKEN = result
	c.write(filename)
}

func (c *Config) configureModel(filename string) {
	if c.MODEL != "" {
		return
	}

	modelList := openai.GetModels()
	var filteredModelList []string
	for _, model := range modelList {
		if strings.Contains(model, "gpt") {
			filteredModelList = append(filteredModelList, model)
		}
	}

	prompt := promptui.Select{
		Label: "🤖 Please select a model",
		Items: filteredModelList,
	}

	_, result, err := prompt.Run()
	if err != nil {
		panic(err)
	}

	c.MODEL = result
	c.write(filename)
}

func Load() error {
	c := &Config{}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	filename := filepath.Join(homeDir, "gpt-quiz", "config.json")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(filename), 0755)
		if err != nil {
			return err
		}
		err = c.write(filename)
		if err != nil {
			return err
		}
	}

	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, c)
	if err != nil {
		return err
	}

	c.configureToken(filename)
	c.configureModel(filename)

	c.setEnvVariables()
	return nil
}

func (c *Config) setEnvVariables() {
	err := os.Setenv("TOKEN", c.TOKEN)
	if err != nil {
		panic(err)
	}
}

func Remove() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	filename := filepath.Join(homeDir, "gpt-quiz", "config.json")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil
	}

	err = os.Remove(filename)
	if err != nil {
		return err
	}

	return nil
}

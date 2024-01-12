package helpers

import (
	"errors"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
)

// PromptContent is the struct that let's us have interactive prompts
type PromptContent struct {
	ErrorMsg string
	Label    string
}

func validate(input string) error {
	if len(input) <= 0 {
		return errors.New("validation error length of input must be greater then 0")
	}
	return nil
}

// GeneratePrompt generates our prompt UI template
func GeneratePrompt(pc PromptContent, v func(string) error) promptui.Prompt {
	validate := v

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	return promptui.Prompt{
		Label:     pc.Label,
		Templates: templates,
		Validate:  validate,
	}
}

func PromptInput(p promptui.Prompt) string {
	result, err := p.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return result
}

func PromptGetInput(pc PromptContent) string {
	p := GeneratePrompt(pc, validate)
	return PromptInput(p)
}

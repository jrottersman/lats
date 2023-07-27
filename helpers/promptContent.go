package helpers

import (
	"errors"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
)

type PromptContent struct {
	ErrorMsg string
	Label    string
}

func GeneratePrompt(pc PromptContent) promptui.Prompt {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(pc.ErrorMsg)
		}
		return nil
	}

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
	p := GeneratePrompt(pc)
	return PromptInput(p)
}

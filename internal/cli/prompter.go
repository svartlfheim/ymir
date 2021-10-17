package cli

import (
	"github.com/manifoldco/promptui"
)

type Prompter struct{}

func (p *Prompter) Ask(q string) (a string, err error) {
	prompt := promptui.Prompt{
		Label: q,
	}

	return prompt.Run()
}

func NewPrompter() *Prompter {
	return &Prompter{}
}

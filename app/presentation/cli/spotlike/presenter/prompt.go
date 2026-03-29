package presenter

import (
	"github.com/yanosea/spotlike/pkg/proxy"
	"github.com/yanosea/spotlike/pkg/utility"
)

var (
	// Pu is a variable that contains the PromptUtil struct for injecting dependencies in testing.
	Pu = utility.NewPromptUtil(proxy.NewPromptui())
)

// RunPrompt runs the prompt.
func RunPrompt(label string) (string, error) {
	prompt := Pu.GetPrompt(label)
	return prompt.Run()
}

// RunPromptWithMask runs the prompt with mask.
func RunPromptWithMask(label string, mask rune) (string, error) {
	prompt := Pu.GetPrompt(label)
	prompt.SetMask(mask)
	return prompt.Run()
}

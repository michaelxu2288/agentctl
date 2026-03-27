package tools

import "fmt"

type CompositeTool struct {
	NameValue        string
	DescriptionValue string
	Steps            []Tool
}

func (c CompositeTool) Name() string        { return c.NameValue }
func (c CompositeTool) Description() string { return c.DescriptionValue }

func (c CompositeTool) Call(params map[string]any) (string, error) {
	if len(c.Steps) == 0 {
		return "", fmt.Errorf("composite tool has no steps")
	}
	output := ""
	for _, step := range c.Steps {
		partial, err := step.Call(params)
		if err != nil {
			return "", err
		}
		if output == "" {
			output = partial
		} else {
			output = output + "\n" + partial
		}
	}
	return output, nil
}

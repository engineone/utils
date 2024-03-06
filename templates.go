package utils

import (
	"bytes"
	"text/template"

	"github.com/engineone/types"
	"github.com/palantir/stacktrace"
)

func RenderInputTemplate(input string, t *types.Task, otherTasks []*types.Task) (string, error) {
	tmpl, err := template.New("input").Parse(input)
	if err != nil {
		return "", stacktrace.PropagateWithCode(err, types.ErrInvalidTask, "Error parsing input template")
	}

	// Create a map of task outputs to pass to the template
	taskOutputs := map[string]interface{}{
		"input": t.GlobalInput,
	}
	for _, id := range t.DependsOn {
		var task *types.Task
		for _, ot := range otherTasks {
			if ot.ID == id {
				task = ot
				break
			}
		}

		if task == nil {
			return "", stacktrace.NewErrorWithCode(types.ErrInvalidTask, "Task dependency not found: %s", id)
		}
		taskOutputs[id] = task.Output
	}

	// Render the template
	renderedInput := bytes.NewBufferString("")
	err = tmpl.Execute(renderedInput, taskOutputs)
	if err != nil {
		return "", stacktrace.PropagateWithCode(err, types.ErrInvalidTask, "Error rendering input template")
	}
	return renderedInput.String(), nil
}

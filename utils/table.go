package utils

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

var (
	levelMapping = map[string]string{
		"normal":  "✅",
		"warning": "⚠️",
		"error":   "❌",
		"unknown": "❓",
	}
)

type Result struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Describe string `json:"describe"`
	Value    string `json:"value"`
	Level    string `json:"level"` //NORMAL, WARNING, ERROR
	Message  string `json:"message"`
}

func PrintResult(causes []*Result) {
	fmt.Println()
	fmt.Println()
	fmt.Println("\033[0;32m +++++++++++++++++++++++++    RUNNING RESULT    +++++++++++++++++++++++++ \033[0m")

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Name", "Description", "Target(s)", "Level", "Message"})

	if len(causes) == 0 {
		msg := "Unable to get any results."
		t.AppendRow(table.Row{"1", "error", msg, "null", levelMapping["unknown"], msg}, table.RowConfig{AutoMerge: true, AutoMergeAlign: text.AlignCenter})
		t.Render()
		return
	}

	for _, cause := range causes {
		t.AppendRow(table.Row{cause.ID, cause.Name, cause.Describe, cause.Value, levelMapping[cause.Level], cause.Message})
		t.AppendSeparator()
	}
	t.Render()
}

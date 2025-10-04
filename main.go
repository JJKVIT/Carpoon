package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

func start(jsonData *jsonLoad, state string) {
	projects := jsonData.Projects
	settings := jsonData.Settings

	m := NewModel(projects, &settings, jsonData, state)

	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		log.Fatal("unable to run")
	}

	m, ok := finalModel.(model)

	if !ok {
		log.Fatal("could not cast model")
	}
	if m.moveFlag {
		if m.moveFlag && len(m.projects) > 0 {
			fmt.Fprintln(os.Stderr, m.projects[m.curr].Path)
		}
	}
}

func main() {
	args := os.Args
	jsonData := &jsonLoad{}
	jsonData.Init()
	if len(args) == 1 {
		start(jsonData, "")
	}

	for _, value := range args {
		if value == "a" {
			dir, err := os.Getwd()
			if err != nil {
				log.Fatal("CANNOT ADD DIRECTORY")
			}
			jsonData.addProject(dir, filepath.Base(dir))
		}
		if value == "e" {
			start(jsonData, "e")
		}
	}

}

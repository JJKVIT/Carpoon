package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	entryView uint = iota
	helpView
	listView
	settingsView
	colorPicker
)

type updateMsg struct{}

type model struct {
	state     uint
	projects  []Project
	config    *Config
	dataobj   *jsonLoad
	width     int
	height    int
	curr      int
	moveFlag  bool
	search    string
	prevState uint
}

func NewModel(projects []Project, config *Config, obj *jsonLoad, state string) model {
	if state == "e" || state == "E" {
		return model{
			state:     listView,
			config:    config,
			projects:  projects,
			dataobj:   obj,
			curr:      0,
			moveFlag:  false,
			search:    "",
			prevState: listView,
		}
	}
	return model{
		state:     entryView,
		config:    config,
		projects:  projects,
		dataobj:   obj,
		curr:      0,
		moveFlag:  false,
		search:    "",
		prevState: entryView,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func addUpdate(m model) []Project {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("CANNOT ADD DIRECTORY")
	}
	m.dataobj.addProject(dir, filepath.Base(dir))
	return m.dataobj.Projects
}

func removeUpdate(m model) []Project {
	m.dataobj.removeProject(m.curr)
	return m.dataobj.Projects
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.prevState != m.state {
		m.search = ""
		m.prevState = m.state
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		key := msg.String()
		if key == "Q" || key == "q" || key == "ctrl+c" {
			return m, tea.Quit
		}
		if key == "A" || key == "a" {
			m.projects = addUpdate(m)
		}
		if _, err := strconv.ParseInt(key, 10, 32); err == nil || m.search != "" {
			if key == "enter" {
				newval, err := strconv.ParseInt(m.search, 10, 32)
				if err != nil {
					fmt.Println("search failed")
				}
				m.curr = int(newval) - 1
				m.moveFlag = true
				return m, tea.Quit
			}
			if key == "backspace" {
				m.search = m.search[:len(m.search)-1]
				if m.search == "" {
					return m, nil
				}
			}
			m.search += key
			intval, err := strconv.ParseInt(m.search, 10, 32)
			if int(intval) > len(m.projects) {
				m.search = strconv.Itoa(len(m.projects))
			} else if err != nil {
				m.search = m.search[:len(m.search)-len(key)]
			}
		}
		switch m.state {
		case entryView:
			switch key {
			case "h", "H":
				m.state = helpView
			case "e", "E":
				m.curr = 0
				m.state = listView
			case "C", "c":
				m.state = settingsView
			case "p":
				m.curr = 0
				m.state = colorPicker
			}
		case listView:
			switch key {
			case "h", "H":
				m.state = helpView
			case "backspace":
				m.state = entryView
			case "C", "c":
				m.state = settingsView
			case "D", "d":
				m.projects = removeUpdate(m)
				if m.curr == len(m.projects) {
					m.curr = m.curr - 1
				}
				if len(m.projects) == 0 {
					m.curr = 0
				}
			case "J", "j", "down":
				m.curr = (m.curr + 1) % len(m.projects)
			case "up", "K", "k":
				m.curr = m.curr - 1
				if m.curr < 0 {
					m.curr = len(m.projects) - 1
				}
			case "enter":
				m.moveFlag = true
				return m, tea.Quit
			}
		case helpView:
			switch key {
			case "e", "E":
				m.curr = 0
				m.state = listView
			case "backspace":
				m.state = entryView
			case "C", "c":
				m.state = settingsView
			case "D", "d":
				m.projects = removeUpdate(m)
				if m.curr == len(m.projects) {
					m.curr = m.curr - 1
				}
			}
		case settingsView:
			switch key {
			case "h", "H":
				m.state = helpView
			case "e", "E":
				m.curr = 0
				m.state = listView
			case "backspace":
				m.state = entryView
			case "J", "j", "down":
				m.curr = (m.curr + 1) % 2
			case "up", "K", "k":
				m.curr = m.curr - 1
				if m.curr < 0 {
					m.curr = 1
				}
			case "D", "d":
				m.projects = removeUpdate(m)
				if m.curr == len(m.projects) {
					m.curr = m.curr - 1
				}
			}
		case colorPicker:
			switch key {
			case "backspace":
				m.state = entryView
			case "enter":
				m.dataobj.changeColor(m.config.SelectColor, m.config.H, m.config.S, m.config.L)
			case "J", "j", "down":
				m.curr = (m.curr + 1) % 3
			case "up", "K", "k":
				m.curr = m.curr - 1
				if m.curr < 0 {
					m.curr = 2
				}
			case "right":
				m.config.SelectColor = HSLtoHEX(m.config.H, m.config.S, m.config.L)
				switch m.curr {
				case 0:
					m.config.H++
					if m.config.H > 360 {
						m.config.H = 360
					}
				case 1:
					m.config.S++
					if m.config.S > 100 {
						m.config.S = 100
					}
				case 2:
					m.config.L++
					if m.config.L > 100 {
						m.config.L = 100
					}
				}
			case "left":
				m.config.SelectColor = HSLtoHEX(m.config.H, m.config.S, m.config.L)
				switch m.curr {
				case 0:
					m.config.H--
					if m.config.H < 0 {
						m.config.H = 0
					}
				case 1:
					m.config.S--
					if m.config.S < 0 {
						m.config.S = 0
					}
				case 2:
					m.config.L--
					if m.config.L < 0 {
						m.config.L = 0
					}
				}
			}
		}
	}
	return m, nil
}

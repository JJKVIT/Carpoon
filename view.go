package main

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	var (
		faintStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color(99)).Faint(true).Align(lipgloss.Center)
		GREEENstlye = lipgloss.NewStyle().Foreground(lipgloss.Color("#4D9740"))
		selectStyle = lipgloss.NewStyle().Background(lipgloss.Color(m.config.SelectColor)).Align(lipgloss.Center)
	)
	content := ""
	controls := ""
	curr := m.curr
	switch m.state {
	case listView:
		for i := range m.projects {
			line := strconv.Itoa(i+1) + " " + m.projects[i].Title + "\n" + m.projects[i].Path + "\n"

			if i == curr {
				line = selectStyle.Render(line)
			}

			content += line + "\n"
		}

		controls = "Q/q (QUIT)  H/h (HELP)  C/c (CONFIG)  ⌫ |Backsapce (HOME)"
	case helpView:
		//IF U KNOW ME DONT ASK WHY I DID IT LIKE THIS... NO ONES FUCKING USING THIS SHIT ANYWAYS
		otherInfo := "\nOther Info\nThe initial list size is set to 10 but it can be increased up to 255 in config\n The select color in list view can be changed in config as well\n\n"
		content = "This is an cli app that is trying to do what harpoon does for " + GREEENstlye.Render("neovim\n\n") + "BASIC COMMANDS:\n" + "1. carpoon - opens the carpoon cli\n" + "2. carpoon a - adds the current directory to the list\n" + "3. carpoon e - opens carpoon list\n" + "4. carpoon (n) - opens the nth directory\n" + "5. d - inside list view d will delete the \n" + "6. a - inside list view a will add the current directory\n" + otherInfo
		controls = "Q/q (QUIT)  E/e(LIST)  C/c (CONFIG) ⌫ |Backsapce (HOME)"
	case entryView:
		content = `█████████████████████████████████████████
 ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░     ██████░
                            ███████░
                    █████░
               ████░

               `
		controls = "Q/q (QUIT)  H/h (HELP)  E/e(LIST)  C/c (CONFIG)"
	case settingsView:
		content = "Max Length of carpoon list : " + strconv.Itoa(m.config.MaxLen) + "\n\n" + "Select Color : " + selectStyle.Render(m.config.SelectColor) + "\n\n"
		controls = "Q/q (QUIT)  E/e(LIST)  H/h (HELP) ⌫ |Backsapce (HOME)"
	}
	if m.width == 0 {
		return "Loading"
	}
	if m.search != "" {
		intval, err := strconv.ParseInt(m.search, 10, 32)
		if err != nil {
			fmt.Println(err)
		}
		controls += "\n\n\n" + m.search + " : " + m.projects[intval-1].Title
	}
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			"CARPOON\n",
			content,
			faintStyle.Width(m.width).Render(controls),
		),
	)
}

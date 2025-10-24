package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func drawGradientBar(value string, m model) string {
	var b string
	for i := range 100 {
		percent := float64(i) / float64(100-1)
		hexColor := ""
		switch value {
		case "h":
			temp := (percent * 360)
			hexColor = HSLtoHEX(temp, m.config.S, m.config.L)
		case "s":
			temp := (percent * 100)
			hexColor = HSLtoHEX(m.config.H, temp, m.config.L)
		case "l":
			temp := percent * 100
			hexColor = HSLtoHEX(m.config.H, m.config.S, temp)
		}

		style := lipgloss.NewStyle().
			Background(lipgloss.Color(hexColor))

		b += (style.Render(" "))
	}

	return b
}

func (m model) View() string {
	if m.width == 0 {
		return "Loading"
	}

	var (
		titleStyle = lipgloss.NewStyle().
				Width(m.width).
				Align(lipgloss.Center).
				MarginBottom(1).
				Bold(true)
		contentStyle = lipgloss.NewStyle().
				Width(m.width).
				Align(lipgloss.Center)
		faintStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("grey")).
				Faint(true).
				Align(lipgloss.Center)
		GREEENstlye = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#4D9740"))
		selectedItemStyle = lipgloss.NewStyle().
					Background(lipgloss.Color(m.config.SelectColor)).
					Align(lipgloss.Center)
		normalItemStyle = lipgloss.NewStyle()
	)

	var content, controls string

	switch m.state {
	case listView:
		var items []string
		for i, project := range m.projects {
			itemStr := fmt.Sprintf("%d. %s\n%s", i+1, project.Title, project.Path)
			if i == m.curr {
				items = append(items, selectedItemStyle.Render(itemStr))
			} else {
				items = append(items, normalItemStyle.Render(itemStr))
			}
		}
		content = strings.Join(items, "\n\n")
		controls = "Q/q (QUIT)  H/h (HELP)  C/c (CONFIG)  ⌫ |Backsapce (HOME)"

	case helpView:
		header := "This is an cli app that is trying to do what harpoon does for " + GREEENstlye.Render("neovim\n\n")
		helpBody := `BASIC COMMANDS:
1. carpoon - opens the carpoon cli
2. carpoon a - adds the current directory to the list
3. carpoon e - opens carpoon list
4. carpoon (n) - opens the nth directory
5. d - inside list view d will delete the
6. a - inside list view a will add the current directory

Other Info
The initial list size is set to 10 but it can be increased up to 255 in config
The select color in list view can be changed in config as well`
		content = header + helpBody
		controls = "Q/q (QUIT)  E/e(LIST)  C/c (CONFIG) ⌫ |Backsapce (HOME)"

	case entryView:
		content = `█████████████████████████████████████████
 ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░     ██████░
                            ███████░
                    █████░
               ████░

               `
		controls = "Q/q (QUIT)  H/h (HELP)  E/e(LIST)  C/c (CONFIG)"

	case settingsView:
		colorBlock := lipgloss.NewStyle().Background(lipgloss.Color(m.config.SelectColor)).Render("   ")
		content = fmt.Sprintf("Max Length of carpoon list : %d\n\nSelect Color : %s", m.config.MaxLen, colorBlock)
		controls = "Q/q (QUIT)  E/e(LIST)  H/h (HELP) ⌫ |Backsapce (HOME)"

	case colorPicker:
		color := m.config.SelectColor
		colorPreviewStyle := lipgloss.NewStyle().
			Padding(2, 5).
			Background(lipgloss.Color(color))

		previewText := lipgloss.JoinVertical(
			lipgloss.Center,
			"Color Preview",
		)

		selectedBox := colorPreviewStyle.Render(previewText)

		hue := drawGradientBar("h", m)
		saturation := drawGradientBar("s", m)
		lightness := drawGradientBar("l", m)

		switch m.curr {
		case 0:
			hue = "> " + hue
		case 1:
			saturation = "> " + saturation
		case 2:
			lightness = "> " + lightness
		}

		content = lipgloss.JoinVertical(
			lipgloss.Center,
			selectedBox+"\n",
			color+"\n",
			fmt.Sprintf("Hue:%d Saturation:%d Lightness:%d\n", int(m.config.H), int(m.config.S), int(m.config.L)),
			hue+"\n",
			saturation+"\n",
			lightness+"\n",
			faintStyle.Render("← → To increase or decrease the selected value\n"),
			faintStyle.Render("↑ ↓ or scroll To change the selected value"),
		)

		controls = "Q/q (QUIT) E/e(LIST) H/h (HELP) ⌫ |Backsapce (HOME)"

	}

	if m.search != "" {
		intval, err := strconv.ParseInt(m.search, 10, 32)
		if err == nil && intval > 0 && int(intval) <= len(m.projects) {
			controls += fmt.Sprintf("\n\n\n%s : %s", m.search, m.projects[intval-1].Title)
		}
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			titleStyle.Render("CARPOON"),
			contentStyle.Render(content),
			faintStyle.Width(m.width).Render("\n\n"+controls),
		),
	)
}

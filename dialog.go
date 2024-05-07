
package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
    "net/http"
    "fmt"
)

var (
	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0)

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			Padding(0, 3).
			MarginTop(1).
			MarginRight(2)

	activeButtonStyle = buttonStyle.Copy().
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94")).
				MarginRight(2).
				Underline(true)
)

type dialog struct {
	id     string
	height int
	width  int

	active   string
	question string

}

func (m dialog) Init() tea.Cmd {
	return nil
}

func (m dialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case tea.MouseMsg:
		if msg.Type != tea.MouseLeft {
			return m, nil
		}
		if zone.Get(m.id + "⏵").InBounds(msg) {
            _,err := http.Get("http://localhost:3000/devices/player/play")
            if err != nil {
                fmt.Println(err)
            }
			m.active = "⏸"
		} else if zone.Get(m.id + "⏸").InBounds(msg) {
            _,err := http.Get("http://localhost:3000/devices/player/pause")
            if err != nil {
                fmt.Println(err)
            }
			m.active = "⏸"
		} else if zone.Get(m.id + "⏮").InBounds(msg) {
            _,err := http.Get("http://localhost:3000/devices/player/previous")
            if err != nil {
                fmt.Println(err)
            }
			m.active = "⏮"
		} else if zone.Get(m.id + "⏭").InBounds(msg){
            _,err := http.Get("http://localhost:3000/devices/player/next")
            if err != nil {
                fmt.Println(err)
            }
            m.active = "⏭"
        }

		return m, nil
	}
	return m, nil
}

func (m dialog) View() string {
    var playButton,playpauseButton,backButton,forwardButton string
    if m.active == "playButton"{
        backButton = buttonStyle.Render("⏮")
        playButton = activeButtonStyle.Render("⏵")
        playpauseButton = activeButtonStyle.Render("⏸")
        forwardButton = buttonStyle.Render("⏭")
    }
    if m.active == "playpauseButton"{
        backButton = buttonStyle.Render("⏮")
        playButton = buttonStyle.Render("⏵")
        playpauseButton = activeButtonStyle.Render("⏸")
        forwardButton = buttonStyle.Render("⏭")
    }
    if m.active == "backButton"{
        backButton = activeButtonStyle.Render("⏮")
        playButton = buttonStyle.Render("⏵")
        playpauseButton = buttonStyle.Render("⏸")
        forwardButton = buttonStyle.Render("⏭")
    }
    if m.active == "forwardButton"{
        backButton = buttonStyle.Render("⏮")
        playButton = buttonStyle.Render("⏵")
        playpauseButton = buttonStyle.Render("⏸")
        forwardButton = activeButtonStyle.Render("⏭")
    } else {
        backButton = buttonStyle.Render("⏮")
        playButton = buttonStyle.Render("⏵")
        playpauseButton = buttonStyle.Render("⏸")
        forwardButton = buttonStyle.Render("⏭")
	}

	question := lipgloss.NewStyle().Width(27).Align(lipgloss.Center).Render("gootify")
	buttons := lipgloss.JoinHorizontal(
		lipgloss.Top,
        zone.Mark(m.id+"⏮",backButton),
        zone.Mark(m.id+"⏵",playButton),
        zone.Mark(m.id+"⏸",playpauseButton),
        zone.Mark(m.id+"⏭",forwardButton),

	)
	return dialogBoxStyle.Render(lipgloss.JoinVertical(lipgloss.Center, question, buttons))
}

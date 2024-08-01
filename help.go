package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type help struct {
	id         string
	height     int
	width      int
	Pages      [][]string
	activePage int
}

func(m help) Init() tea.Cmd{
    return nil
}

func (m help) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case ".":
			m.activePage = min(m.activePage+1, len(m.Pages)-1)
			return m, nil
		case ",":
			m.activePage = max(m.activePage-1, 0)
			return m, nil
		}
	}
	return m, nil
}

func (m help) View() string {
	for i := range m.Pages {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.Pages)-1, i == m.activePage

		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
	}

    styles := []string{""}
    for _,x := range m.Pages[m.activePage]{
        styles = append(styles,lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(x))
    }
    return docStyle.Render(lipgloss.JoinVertical(lipgloss.Center,styles[0:]...))
}

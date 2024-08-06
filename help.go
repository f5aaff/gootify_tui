package main

import (
    "fmt"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

type CurrentPage struct {
    lines []string
    page  int
}

var currentPage = CurrentPage{lines: []string{}, page: 0}

type help struct {
    id         string
    height     int
    width      int
    Pages      [][]string
    activePage int
}

func (m help) Init() tea.Cmd {
    return nil
}

func (m help) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch keypress := msg.String(); keypress {
        case "ctrl+c", "q":
            return m, tea.Quit
        case "]":
            m.activePage = min(m.activePage+1, len(m.Pages)-1)
            currentPage.page = m.activePage
            currentPage.lines = m.Pages[m.activePage]
        case "[":
            m.activePage = max(m.activePage-1, 0)
            currentPage.page = m.activePage
            currentPage.lines = m.Pages[m.activePage]
        }
    }
    return m, nil
}

func (m help) View() string {

    makeLine := func(linearr []string) []string {
        var res = []string{}
        for x := range linearr {
            line := lipgloss.NewStyle().Width(m.width).Margin(0,10,1).Align(lipgloss.Left).Render(linearr[x])
            res = append(res, line)
        }
        return res
    }

    styles := makeLine(currentPage.lines)

    styles = append(styles, lipgloss.NewStyle().Width(m.width).Margin(0,10,1).Align(lipgloss.Center).Render(fmt.Sprint(currentPage.page)))
    page := lipgloss.JoinVertical(lipgloss.Center,styles...)
    return docStyle.Render(page)
}

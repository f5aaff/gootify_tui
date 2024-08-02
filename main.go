package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

const (
    baseURL = "http://localhost:3000/"
)

var (
    subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
    highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
)

func tabBorderWithBottom(left, middle, right, top string) lipgloss.Border {
    border := lipgloss.RoundedBorder()
    border.BottomLeft = left
    border.Bottom = middle
    border.BottomRight = right
    border.Top = top
    return border
}

var (
    inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴", "─")
    activeTabBorder   = tabBorderWithBottom("┘", "_", "└", "─")
    docStyle          = lipgloss.NewStyle().Padding(1, 1, 1, 1)
    highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
    inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
    activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
    windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(0, 0).Align(lipgloss.Center).Border(lipgloss.RoundedBorder())
)

func main() {

    go UpdateInterval()

    zone.NewGlobal()
    tabs := []string{"Dialog", "Help"}
    d := dialog{width: 20, height: 20}
    playerHelp := []string{"ctrl+c/q : quit", "d : next", "l/a : previous", "s : pause", "shift+s : play","+ : volume up", "- : volume down"}
    helpHelp := []string{"help", "+ : next help page", "- : previous help page"}
    h := help{}
    h.Pages = append(h.Pages,playerHelp)
    h.Pages = append(h.Pages,helpHelp)
    h.width = d.width
    h.height = d.height
    tabContent := []subModel{d, h}
    m := mainModel{Tabs: tabs, TabContent: tabContent}
    if _, err := tea.NewProgram(m).Run(); err != nil {
        fmt.Println("Error running program:", err)
        os.Exit(1)
    }
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

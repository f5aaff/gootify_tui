// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

// This is a modified version of this example, supporting full screen, dynamic
// resizing, and clickable models (tabs, lists, dialogs, etc).
// 	https://github.com/charmbracelet/lipgloss/blob/master/example

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	divider = lipgloss.NewStyle().
		SetString("•").
		Padding(0, 1).
		Foreground(subtle).
		String()
)

type model struct {
	height int
	width  int

	dialog  tea.Model

}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) isInitialized() bool {
	return m.height != 0 && m.width != 0
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.isInitialized() {
		if _, ok := msg.(tea.WindowSizeMsg); !ok {
			return m, nil
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Example of toggling mouse event tracking on/off.
		if msg.String() == "ctrl+e" {
			zone.SetEnabled(!zone.Enabled())
			return m, nil
		}

		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		msg.Height -= 30
		msg.Width -= 30
		return m.propagate(msg), nil
	}

	return m.propagate(msg), nil
}

func (m *model) propagate(msg tea.Msg) tea.Model {
	// Propagate to all children.
	m.dialog, _ = m.dialog.Update(msg)
return m
}

func (m model) View() string {
	if !m.isInitialized() {
		return ""
	}

	s := lipgloss.NewStyle().MaxHeight(m.height).MaxWidth(m.width).Padding(1, 2, 1, 2)

	return zone.Scan(s.Render(lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.PlaceHorizontal(
			m.width, lipgloss.Center,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
                m.dialog.View(),
			),
			lipgloss.WithWhitespaceChars(" "),
		),
	)))
}

func main() {
	// Initialize a global zone manager, so we don't have to pass around the manager
	// throughout components.
	zone.NewGlobal()

	m := &model{
		dialog: &dialog{
			id:       zone.NewPrefix(),
			height:   50,
            width:    50,
			active:   "playpauseButton",
		},

	}

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

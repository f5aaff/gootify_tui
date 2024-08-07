// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"gootifyTui/models"
	"io"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var currentItem int = 0
var (
	listStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(subtle).
			Margin(0, 0, 1)

	listStyleCurrent = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, true, false).
				BorderForeground(highlight).
				Margin(0, 0, 1)
	listHeader = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(subtle).
			Margin(0, 10, 1).
			Render

	listItemStyle = lipgloss.NewStyle().PaddingLeft(2).Render

	checkMark = lipgloss.NewStyle().SetString("✓").
			Foreground(highlightColor).
			PaddingRight(1).
			String()

	createListItem = func(s string, b bool) string {
		if b {
			return listStyleCurrent.Render(s)
		}
		return listStyle.Render(s)
	}
	listDoneStyle = func(s string) string {
		return checkMark + lipgloss.NewStyle().
			Strikethrough(true).
			Foreground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
			Render(s)
	}
)

type listItem struct {
	name string
	done bool
}

type list struct {
	id     string
	height int
	width  int

	title      string
	items      []listItem
	activeItem int
}

func (m list) Init() tea.Cmd {
	return nil
}

func (m list) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "s", "j":
			m.activeItem = min(m.activeItem+1, len(m.items)-1)
			currentItem = m.activeItem
			fmt.Print(m.activeItem)
			return m, nil
		case "w", "k":
			m.activeItem = max(m.activeItem-1, 0)
			currentItem = m.activeItem
			fmt.Print(m.activeItem)
			return m, nil
		case "return", "d":
			if m.items[currentItem].done {
				m.items[currentItem].done = false
				return m, nil
			}
			m.items[currentItem].done = true
			return m, nil
		}
	}

	return m, nil
}

func (m list) View() string {
	out := []string{listHeader(m.title)}

	for x, item := range m.items {
		selected := false
		if currentItem == x {
			selected = true
		}

		item.name = createListItem(item.name, selected)
		if item.done {
			out = append(out, checkMark+item.name)
			continue
		}

		out = append(out, listItemStyle(item.name))
	}

	return listStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, out...),
	)
}

func getQueue(queue *models.Queue) error {
	res, err := http.Get(baseURL + "devices/queue")
	if err != nil {
		return err
	}
	byteres, err := io.ReadAll(res.Body)
	if err != nil{
		return err
	}
	err = json.Unmarshal(byteres, &queue)
	if err != nil{
		return err
	}
	return nil
}

// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"encoding/json"
	"example/models"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"net/http"
)

func getCurrentlyPlaying() string {
	res, err := http.Get(baseURL + "devices/currently_playing")
	if err != nil {
		return err.Error()
	}

	return getTitleAndArtist(res)
}

func getTitleAndArtist(res *http.Response) string {
	var results models.Currently_playing
	byteres, err := io.ReadAll(res.Body)

	if err != nil {
		return err.Error()
	}
	err = json.Unmarshal(byteres, &results)
	if err != nil {
		return err.Error()
	}
	if results.Item.Type == "track" {
		track := results.Item.Name
		artist := results.Item.Artists[0].Name
		return fmt.Sprintf("%s - %s", track, artist)

	}
	return "not a song"
}

var (
	dialogBoxStyle = lipgloss.NewStyle().
			BorderForeground(highlightColor).
			Padding(2, 0).
			Align(lipgloss.Center).
			Border(lipgloss.NormalBorder()).UnsetBorderTop()
	//dialogBoxStyle = lipgloss.NewStyle().
	//		Border(lipgloss.RoundedBorder(), true).
	//		BorderForeground(lipgloss.Color("#874BFD")).
	//		Padding(1, 0)

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
        case tea.KeyMsg:
            switch keypress := msg.String(); keypress {
            case "ctrl+c","q":
                return m, tea.Quit
            }
	}
	return m, nil

}

func (m dialog) View() string {

	question := lipgloss.NewStyle().Width(30).Align(lipgloss.Center).Render("gootify")

	volLabel := lipgloss.NewStyle().Width(30).Align(lipgloss.Center).Render("volume")
	volumeControls := lipgloss.JoinHorizontal(lipgloss.Bottom,
		volLabel)

	currentTrack := lipgloss.NewStyle().Width(30).Align(lipgloss.Center).Render(getCurrentlyPlaying())
	return dialogBoxStyle.Render(lipgloss.JoinVertical(lipgloss.Center, question, currentTrack, volumeControls))
}

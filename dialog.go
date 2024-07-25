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
	zone "github.com/lrstanley/bubblezone"
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
			_, err := http.Get(baseURL + "devices/player/play")
			if err != nil {
				fmt.Println(err)
			}
			m.active = "⏸"
		} else if zone.Get(m.id + "⏸").InBounds(msg) {
			_, err := http.Get(baseURL + "devices/player/pause")
			if err != nil {
				fmt.Println(err)
			}
			m.active = "⏸"
		} else if zone.Get(m.id + "⏮").InBounds(msg) {
			_, err := http.Get(baseURL + "devices/player/previous")
			if err != nil {
				fmt.Println(err)
			}
			m.active = "⏮"
		} else if zone.Get(m.id + "⏭").InBounds(msg) {
			_, err := http.Get(baseURL + "devices/player/next")
			if err != nil {
				fmt.Println(err)
			}
			m.active = "⏭"
		} else if zone.Get(m.id + "+").InBounds(msg) {
			_, err := http.Get(baseURL + "devices/volup")
			if err != nil {
				fmt.Println(err)
			}
		} else if zone.Get(m.id + "-").InBounds(msg) {
			_, err := http.Get(baseURL + "devices/voldown")
			if err != nil {
				fmt.Println(err)
			}
		}

		return m, nil
	}
	return m, nil

}

func (m dialog) View() string {
	var playButton, playpauseButton, backButton, forwardButton, volUpButton, volDownButton string
	if m.active == "playButton" {
		backButton = buttonStyle.Render("⏮")
		playButton = activeButtonStyle.Render("⏵")
		playpauseButton = activeButtonStyle.Render("⏸")
		forwardButton = buttonStyle.Render("⏭")
		volUpButton = buttonStyle.Render("+")
		volDownButton = buttonStyle.Render("-")
	}
	if m.active == "playpauseButton" {
		backButton = buttonStyle.Render("⏮")
		playButton = buttonStyle.Render("⏵")
		playpauseButton = activeButtonStyle.Render("⏸")
		forwardButton = buttonStyle.Render("⏭")
		volUpButton = buttonStyle.Render("+")
		volDownButton = buttonStyle.Render("-")
	}
	if m.active == "backButton" {
		backButton = activeButtonStyle.Render("⏮")
		playButton = buttonStyle.Render("⏵")
		playpauseButton = buttonStyle.Render("⏸")
		forwardButton = buttonStyle.Render("⏭")
		volUpButton = buttonStyle.Render("+")
		volDownButton = buttonStyle.Render("-")
	}

	if m.active == "forwardButton" {
		backButton = buttonStyle.Render("⏮")
		playButton = buttonStyle.Render("⏵")
		playpauseButton = buttonStyle.Render("⏸")
		forwardButton = activeButtonStyle.Render("⏭")
		volUpButton = buttonStyle.Render("+")
		volDownButton = buttonStyle.Render("-")
	} else {
		backButton = buttonStyle.Render("⏮")
		playButton = buttonStyle.Render("⏵")
		playpauseButton = buttonStyle.Render("⏸")
		forwardButton = buttonStyle.Render("⏭")
		volUpButton = buttonStyle.Render("+")
		volDownButton = buttonStyle.Render("-")
	}

	question := lipgloss.NewStyle().Width(27).Align(lipgloss.Center).Render("gootify")
	buttons := lipgloss.JoinHorizontal(
		lipgloss.Top,
		zone.Mark(m.id+"⏮", backButton),
		zone.Mark(m.id+"⏵", playButton),
		zone.Mark(m.id+"⏸", playpauseButton),
		zone.Mark(m.id+"⏭", forwardButton),
	)

	volLabel := lipgloss.NewStyle().Width(27).Align(lipgloss.Center).Render("volume")
	volumeControls := lipgloss.JoinHorizontal(lipgloss.Bottom,
		volLabel,
		zone.Mark(m.id+"+", volUpButton),
		zone.Mark(m.id+"-", volDownButton),
	)

	currentTrack := lipgloss.NewStyle().Width(27).Align(lipgloss.Center).Render(getCurrentlyPlaying())
	return dialogBoxStyle.Render(lipgloss.JoinVertical(lipgloss.Center, question, buttons, currentTrack, volumeControls))
}

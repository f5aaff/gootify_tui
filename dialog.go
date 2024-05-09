package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
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

type Currently_playing struct {
	Timestamp int64 `json:"timestamp"`
	Context   struct {
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href string `json:"href"`
		Type string `json:"type"`
		URI  string `json:"uri"`
	} `json:"context"`
	ProgressMs int `json:"progress_ms"`
	Item       struct {
		Album struct {
			AlbumType string `json:"album_type"`
			Artists   []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href string `json:"href"`
				ID   string `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
				URI  string `json:"uri"`
			} `json:"artists"`
			AvailableMarkets []string `json:"available_markets"`
			ExternalUrls     struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href   string `json:"href"`
			ID     string `json:"id"`
			Images []struct {
				Height int    `json:"height"`
				URL    string `json:"url"`
				Width  int    `json:"width"`
			} `json:"images"`
			Name                 string `json:"name"`
			ReleaseDate          string `json:"release_date"`
			ReleaseDatePrecision string `json:"release_date_precision"`
			TotalTracks          int    `json:"total_tracks"`
			Type                 string `json:"type"`
			URI                  string `json:"uri"`
		} `json:"album"`
		Artists []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href string `json:"href"`
			ID   string `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
			URI  string `json:"uri"`
		} `json:"artists"`
		AvailableMarkets []string `json:"available_markets"`
		DiscNumber       int      `json:"disc_number"`
		DurationMs       int      `json:"duration_ms"`
		Explicit         bool     `json:"explicit"`
		ExternalIds      struct {
			Isrc string `json:"isrc"`
		} `json:"external_ids"`
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href        string `json:"href"`
		ID          string `json:"id"`
		IsLocal     bool   `json:"is_local"`
		Name        string `json:"name"`
		Popularity  int    `json:"popularity"`
		PreviewURL  string `json:"preview_url"`
		TrackNumber int    `json:"track_number"`
		Type        string `json:"type"`
		URI         string `json:"uri"`
	} `json:"item"`
	CurrentlyPlayingType string `json:"currently_playing_type"`
	Actions              struct {
		Disallows struct {
			Resuming     bool `json:"resuming"`
			SkippingPrev bool `json:"skipping_prev"`
		} `json:"disallows"`
	} `json:"actions"`
	IsPlaying bool `json:"is_playing"`
}

func getCurrentlyPlaying() string {
	res, err := http.Get("http://localhost:3000/devices/currently_playing")
	if err != nil {
		return err.Error()
	}

	return getTitleAndArtist(res)
}

func getTitleAndArtist(res *http.Response) string {
	var results Currently_playing
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
			_, err := http.Get("http://localhost:3000/devices/player/play")
			if err != nil {
				fmt.Println(err)
			}
			m.active = "⏸"
		} else if zone.Get(m.id + "⏸").InBounds(msg) {
			_, err := http.Get("http://localhost:3000/devices/player/pause")
			if err != nil {
				fmt.Println(err)
			}
			m.active = "⏸"
		} else if zone.Get(m.id + "⏮").InBounds(msg) {
			_, err := http.Get("http://localhost:3000/devices/player/previous")
			if err != nil {
				fmt.Println(err)
			}
			m.active = "⏮"
		} else if zone.Get(m.id + "⏭").InBounds(msg) {
			_, err := http.Get("http://localhost:3000/devices/player/next")
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
	var playButton, playpauseButton, backButton, forwardButton string
	if m.active == "playButton" {
		backButton = buttonStyle.Render("⏮")
		playButton = activeButtonStyle.Render("⏵")
		playpauseButton = activeButtonStyle.Render("⏸")
		forwardButton = buttonStyle.Render("⏭")
	}
	if m.active == "playpauseButton" {
		backButton = buttonStyle.Render("⏮")
		playButton = buttonStyle.Render("⏵")
		playpauseButton = activeButtonStyle.Render("⏸")
		forwardButton = buttonStyle.Render("⏭")
	}
	if m.active == "backButton" {
		backButton = activeButtonStyle.Render("⏮")
		playButton = buttonStyle.Render("⏵")
		playpauseButton = buttonStyle.Render("⏸")
		forwardButton = buttonStyle.Render("⏭")
	}

	if m.active == "forwardButton" {
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
		zone.Mark(m.id+"⏮", backButton),
		zone.Mark(m.id+"⏵", playButton),
		zone.Mark(m.id+"⏸", playpauseButton),
		zone.Mark(m.id+"⏭", forwardButton),
	)
	currentTrack := lipgloss.NewStyle().Width(27).Align(lipgloss.Center).Render(getCurrentlyPlaying())
	return dialogBoxStyle.Render(lipgloss.JoinVertical(lipgloss.Center, question, buttons, currentTrack))
}

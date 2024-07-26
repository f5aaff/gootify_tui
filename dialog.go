package main

import (
	"encoding/json"
	"example/models"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
)

func getCurrentlyPlaying() string {
	res, err := http.Get(baseURL + "devices/currently_playing")
	if err != nil {
		return err.Error()
	}

	return getTitleAndArtist(res)
}
func getVolume() string {
    res,err := http.Get(baseURL + "devices/")
    if err != nil{
        return err.Error()
    }

    byteres,err := io.ReadAll(res.Body)
    if err !=nil{
        return err.Error()
    }

    var device models.DeviceResponse
    err = json.Unmarshal(byteres, &device)
    if err != nil {
        return err.Error()
    }
    return fmt.Sprint(device.Device.VolumePercent)
}
func renderVolume() string {
    volstr := getVolume()
    vol,err := strconv.Atoi(volstr)
    if err != nil {
        return err.Error()
    }
    currentBlock := vol/10
    emptyblocks := 10 - vol
    var s strings.Builder
    for range currentBlock{
        s.WriteString("█")
    }
    for range emptyblocks{
        s.WriteString("-")
    }
    volout := "vol:"+s.String()+" "+fmt.Sprint(vol)+"%"
    return volout
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
)

type dialog struct {
	id     string
	height int
	width  int

	active   string
	question string
}

func playerFunc(action string) bool {
	res, err := http.Get(baseURL + "devices/player/" + action)
	if err != nil {
		getCurrentlyPlaying()
		log.Error().Msg("oops")
	}
	if res.StatusCode == 200 {
		return true
	} else {
		return false
	}
}

func (m dialog) Init() tea.Cmd {
	return nil
}

func (m dialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "n", "d":
			playerFunc("next")
		case "l", "a":
			playerFunc("previous")
		case "s":
			playerFunc("pause")
        case "S":
            playerFunc("play")
		}
	}
	return m, nil

}

func (m dialog) View() string {

	question := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render("gootify")

	volLabel := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(renderVolume())
	volumeControls := lipgloss.JoinHorizontal(lipgloss.Bottom,
		volLabel)
	lipgloss.JoinHorizontal(lipgloss.Center, "\n")
	currentTrack := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(getCurrentlyPlaying())
	return dialogBoxStyle.Render(lipgloss.JoinVertical(lipgloss.Center, question, currentTrack, volumeControls))
}

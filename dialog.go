package main

import (
	"encoding/json"
	"example/models"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

type Current struct {
	track    string
	progress string
	vol      int
}

var current = Current{track: "", progress: "", vol: 0}
var playing = Current{track: "", progress: "", vol: 0}

func saveCover(url string, id string, playing Current) chan string {
	if playing.track != current.track {
		//TODO: make the retrieval a background process, and add a check to see if
		// the cover is already downloaded.
		r := make(chan string)
		// Execute the first command: wget
		go func() {
			path := fmt.Sprintf("albumArt/%s.jpg", id)
			wgetCmd := exec.Command("wget", url, "-O", path)
			if err := wgetCmd.Run(); err != nil {
				r <- "" //, fmt.Errorf("error executing wget command: %v", err)
			}

			// Execute the second command: viu
			viuCmd := exec.Command("viu", "-w", "30", "-h", "10", path)
			viuOutput, err := viuCmd.Output()
			if err != nil {
				r <- "" // fmt.Errorf("error executing viu command: %v,%s", err,url)
			}

			r <- string(viuOutput)
		}()
		return r
	}
	return nil
}

func getAlbumCover(playing Current) string {
	res, err := http.Get(baseURL + "devices/currently_playing")
	if err != nil {
		return err.Error()
	}
	var results models.Currently_playing
	byteres, err := io.ReadAll(res.Body)
	if err != nil {
		return err.Error()
	}
	err = json.Unmarshal(byteres, &results)
	if err != nil {
		return err.Error()
	}
	if results.CurrentlyPlayingType != "episode" {
		url := &results.Item.Album.Images[0].URL
		image := saveCover(*url, results.Item.ID, playing)
		return <-image
	}
	return ""
}
func renderAlbumCover(path string, playing Current) string {
	if playing.track != current.track {
		viuCmd := exec.Command("viu", "-w", "30", "-h", "10", path)
		viuOutput, err := viuCmd.Output()
		if err != nil {
			return ""
		}
		return string(viuOutput)
	}
	return ""
}

func getCurrentlyPlaying() string {
	res, err := http.Get(baseURL + "devices/currently_playing")
	if err != nil {
		return err.Error()
	}

	return getTitleAndArtist(res)
}
func getVolume() string {
	res, err := http.Get(baseURL + "devices/")
	if err != nil {
		return err.Error()
	}

	byteres, err := io.ReadAll(res.Body)
	if err != nil {
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
	vol, err := strconv.Atoi(volstr)
	if err != nil {
		return err.Error()
	}
	currentBlock := vol / 10
	emptyblocks := 10 - vol
	var s strings.Builder
	for range currentBlock {
		s.WriteString("█")
	}
	for range emptyblocks {
		s.WriteString("-")
	}
	volout := "vol:" + s.String() + " " + fmt.Sprint(vol) + "%"
	return volout
}

func getTitleAndArtist(res *http.Response) string {
	if playing.track != current.track {

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
			current.track = fmt.Sprintf("%s - %s", track, artist)
			return current.track

		}
		if results.CurrentlyPlayingType == "episode" {
			track := results.Item.Name
			return track
		}
		return "not a song"
	}
	return current.track
}

var (
	dialogBoxStyle = lipgloss.NewStyle().
		BorderForeground(highlightColor).
		Padding(2, 0).
		Align(lipgloss.Center).
		Border(lipgloss.RoundedBorder())
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
	playing.track = getCurrentlyPlaying()
	volLabel := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(renderVolume())
	volumeControls := lipgloss.JoinHorizontal(lipgloss.Bottom, volLabel)
	currentTrack := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(current.track)
	//albumArt := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(renderAlbumCover("./alb.jpg", playing))
	//fmt.Println(renderAlbumCover("./alb.jpg", playing))
	return dialogBoxStyle.Render(lipgloss.JoinVertical(lipgloss.Center, question, currentTrack, volumeControls))
}

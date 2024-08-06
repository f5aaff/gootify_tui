package main

import (
	"encoding/json"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
	"gootifyTui/models"
	"io"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Current struct {
	track          string
	progress       int
	duration       int
	progressString string
	vol            string
	albumURL       string
	album          string
	currentType    string
	itemID         string
}

var current = Current{track: "", progress: 0, duration: 0, progressString: "", vol: "", albumURL: "", album: "", currentType: "", itemID: ""}

func UpdateInterval(now bool) {
	if now {
		_ = getCurrentlyPlaying()
		renderProgress()
		_ = renderVolume()
		_ = getAlbumCover()
	}
	for range time.Tick(time.Second * 5) {
		_ = getCurrentlyPlaying()
		renderProgress()
		_ = renderVolume()
		_ = getAlbumCover()

	}
}

func saveCover(url string, id string) error {
	// Execute the first command: wget
	path := fmt.Sprintf("albumArt/%s.jpg", id)
	wgetCmd := exec.Command("wget", url, "-O", path)
	if err := wgetCmd.Run(); err != nil {
		return err
	}
	// Execute the second command: viu
	viuCmd := exec.Command("viu", "-b", "-w", "20", "-h", "10", path)
	viuOutput, err := viuCmd.Output()
	if err != nil {
		return err
	}
	if current.album != string(viuOutput) {
		current.album = string(viuOutput)
	}

	return nil
}

func getAlbumCover() error {
	if current.itemID != "" && current.albumURL != "" {
		err := saveCover(current.albumURL, current.itemID)
		if err != nil {
			return err
		}
	}

	return nil
}

func getCurrentlyPlaying() error {
	res, err := http.Get(baseURL + "devices/currently_playing")
	if err != nil {
		return err
	}

	err = getTitleAndArtist(res)
	if err != nil {
		return err
	}
	return nil
}

func getVolume() string {
	res, err := http.Get(baseURL + "devices/")
	if err != nil {
		return err.Error()
	}

	type response struct {
		ID               string `json:"id"`
		IsActive         bool   `json:"is_active"`
		IsPrivateSession bool   `json:"is_private_session"`
		IsRestricted     bool   `json:"is_restricted"`
		Name             string `json:"name"`
		SupportsVolume   bool   `json:"supports_volume"`
		Type             string `json:"type"`
		VolumePercent    int    `json:"volume_percent"`
	}

	byteres, err := io.ReadAll(res.Body)
	if err != nil {
		return err.Error()
	}

	var device response
	err = json.Unmarshal(byteres, &device)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprint(device.VolumePercent)
}
func renderVolume() error {
	volstr := getVolume()
	vol, err := strconv.Atoi(volstr)
	if err != nil {
		return err
	}
	currentBlock := vol / 10
	emptyblocks := 10 - currentBlock
	var s strings.Builder
	for range currentBlock {
		s.WriteString("█")
	}
	for range emptyblocks {
		s.WriteString("-")
	}
	volout := "vol: " + s.String() + " " + fmt.Sprint(vol) + "%"
	if current.vol != volout {
		current.vol = volout
		return nil
	}
	return nil
}
func renderProgress() {

	prog := current.progress

	dur := current.duration
	perc := 0
	if dur > 0 {
		progf := float64(current.progress)
		durf := float64(current.duration)
		perc = int((progf / durf) * 100)
	}
	currentBlock := (100 - perc) / 10
	emptyblocks := 10 - currentBlock
	var s strings.Builder
	for range currentBlock-1 {
		s.WriteString("█")
	}
	for range emptyblocks-1 {
		s.WriteString("-")
	}
	minSec := func(ms int) string {
		if ms != 0 {
			minutes := int((ms / 1000) / 60)
			remainingSeconds := int((ms / 1000) % 60)
			remstr, minstr := fmt.Sprint(remainingSeconds), fmt.Sprint(minutes)
			if minutes < 10 {
				minstr = "0" + fmt.Sprint(minutes)
			}
			if remainingSeconds < 10 {
				remstr = "0" + fmt.Sprint(remainingSeconds)
			}
			return minstr + ":" + remstr
		}
		return "00:00"
	}

	progString := minSec(prog)
	durString := minSec(dur)
	current.progressString = progString + s.String() + durString
}
func getTitleAndArtist(res *http.Response) error {
	byteres, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	resultMap := map[string]interface{}{}
	err = json.Unmarshal(byteres, &resultMap)
	if err != nil {
		return err
	}
	currentType, ok := resultMap["currently_playing_type"]
	if !ok {
		return nil
	}
	progress, ok := resultMap["progress_ms"]
	if !ok {
		progress = "0"
	}

	//this is gross, I hate it but it works
	current.progress = int(progress.(float64))
	if currentType.(string) == "track" {
		current.currentType = "track"
		var results models.Currently_playing
		err = json.Unmarshal(byteres, &results)
		if err != nil {
			return err
		}
		current.duration = results.Item.DurationMs

		track := results.Item.Name
		artist := results.Item.Artists[0].Name
		trackString := fmt.Sprintf("%s - %s", track, artist)

		if trackString != current.track {
			current.track = trackString
			current.itemID = results.Item.ID
			if len(results.Item.Album.Images) > 0 {
				current.albumURL = results.Item.Album.Images[0].URL
			}
			return nil
		}
		return nil

	}

	if currentType.(string) == "episode" {
		current.currentType = "episode"
		res, err := http.Get(baseURL + "devices/queue")
		if err != nil {
			return err
		}
		byteres, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		var results models.Episode
		err = json.Unmarshal(byteres, &results)
		if err != nil {
			return err
		}
		current.duration = results.CurrentlyPlaying.DurationMs
		track := results.CurrentlyPlaying.Name + " - " + results.CurrentlyPlaying.Show.Name
		if current.track != track {
			current.track = track
			current.itemID = results.CurrentlyPlaying.ID
			if len(results.CurrentlyPlaying.Images) > 0 {
				current.albumURL = results.CurrentlyPlaying.Images[0].URL
			}
			return nil
		}
	}
	return nil
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
		err := getCurrentlyPlaying()
		if err != nil {
			return false
		}
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
		case "+":
			_, _ = http.Get(baseURL + "devices/volup")

		case "-":
			_, _ = http.Get(baseURL + "devices/voldown")
		}

	}
	return m, nil

}

func (m dialog) View() string {

	question := lipgloss.NewStyle().Width(m.width).Margin(1).Align(lipgloss.Center).Render("gootify")

	volLabel := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(current.vol)
	volumeControls := lipgloss.JoinHorizontal(lipgloss.Bottom, volLabel)

	track, artist := "", ""

	if current.track != "" {
		track = strings.Split(current.track, " - ")[0]
		artist = strings.Split(current.track, " - ")[1]
	}

	progress := current.progressString

	currentTrack := lipgloss.NewStyle().Width(m.width).Margin(0, 10).Align(lipgloss.Center).Render(track)
	currentArtist := lipgloss.NewStyle().Width(m.width).Margin(0, 10).Align(lipgloss.Center).Render(artist)
	currentProgress := lipgloss.NewStyle().Width(m.width).Margin(0,10,1).Align(lipgloss.Center).Render(progress)
	currentString := lipgloss.JoinVertical(0.1, currentTrack, currentArtist)
	albumArt := lipgloss.NewStyle().Width(m.width).Margin(0, 10).Align(lipgloss.Center).Render(current.album)
	album := lipgloss.JoinVertical(1, currentString, albumArt, currentProgress)
	return dialogBoxStyle.Render(lipgloss.JoinVertical(lipgloss.Center, question, album, volumeControls))
}

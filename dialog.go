package main

import (
    "encoding/json"
    "gootifyTui/models"
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
    vol      string
    albumURL string
    album    string
}

var current = Current{track: "", progress: "", vol: "", albumURL: "", album: ""}

func saveCover(url string, id string) error {
    // Execute the first command: wget
    path := fmt.Sprintf("albumArt/%s.jpg", id)
    wgetCmd := exec.Command("wget", url, "-O", path)
    if err := wgetCmd.Run(); err != nil {
        return err
    }

    // Execute the second command: viu
    viuCmd := exec.Command("viu", "-w", "20", "-h", "10", path)
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
    res, err := http.Get(baseURL + "devices/currently_playing")
    if err != nil {
        return err
    }
    var results models.Currently_playing
    byteres, err := io.ReadAll(res.Body)
    if err != nil {
        return err
    }
    err = json.Unmarshal(byteres, &results)
    if err != nil {
        return err
    }
    if results.CurrentlyPlayingType != "episode" {
        url := &results.Item.Album.Images[0].URL
        if current.albumURL != *url {
            current.albumURL = *url
            err := saveCover(*url, results.Item.ID)
            if err != nil {
                return err
            }
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
    volout := "vol:" + s.String() + " " + fmt.Sprint(vol) + "%"
    if current.vol != volout {
        current.vol = volout
        return nil
    }
    return nil
}

func getTitleAndArtist(res *http.Response) error {
    var results models.Currently_playing
    byteres, err := io.ReadAll(res.Body)

    if err != nil {
        return err
    }
    err = json.Unmarshal(byteres, &results)
    if err != nil {
        return err
    }
    if results.Item.Type == "track" {
        track := results.Item.Name
        artist := results.Item.Artists[0].Name
        trackString := fmt.Sprintf("%s - %s", track, artist)
        if trackString != current.track {
            current.track = trackString
            return nil
        }
        return nil

    }
    if results.CurrentlyPlayingType == "episode" {
        track := results.Item.Name
        if current.track != track {
            current.track = track
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
            _,_ = http.Get(baseURL + "devices/volup")

        case "-":
            _,_ = http.Get(baseURL + "devices/voldown")
        }

    }
    return m, nil

}

func (m dialog) View() string {

    question := lipgloss.NewStyle().Width(m.width).Margin(1).Align(lipgloss.Center).Render("gootify")

    _ = getCurrentlyPlaying()
    _ = renderVolume()
    _ = getAlbumCover()

    volLabel := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(current.vol)
    volumeControls := lipgloss.JoinHorizontal(lipgloss.Bottom, volLabel)
    track := strings.Split(current.track, " - ")[0]
    artist := strings.Split(current.track, " - ")[1]
    currentTrack := lipgloss.NewStyle().Width(m.width).Margin(0, 10).Align(lipgloss.Center).Render(track)
    currentArtist := lipgloss.NewStyle().Width(m.width).Margin(0, 10).Align(lipgloss.Center).Render(artist)

    currentString := lipgloss.JoinVertical(0.1, currentTrack, currentArtist)
    //currentTrack := lipgloss.NewStyle().Width(m.width).Margin(0,10).Align(lipgloss.Center).Render(current.track)
    albumArt := lipgloss.NewStyle().Width(m.width).Margin(0, 10).Align(lipgloss.Center).Render(current.album)
    album := lipgloss.JoinVertical(0.3, currentString, albumArt)
    return dialogBoxStyle.Render(lipgloss.JoinVertical(lipgloss.Center, question, album, volumeControls))
}

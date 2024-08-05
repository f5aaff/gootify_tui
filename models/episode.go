package models

type Episode struct {
	CurrentlyPlaying struct {
		AudioPreviewURL string `json:"audio_preview_url"`
		Description     string `json:"description"`
		DurationMs      int    `json:"duration_ms"`
		Explicit        bool   `json:"explicit"`
		ExternalUrls    struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href            string `json:"href"`
		HTMLDescription string `json:"html_description"`
		ID              string `json:"id"`
		Images          []struct {
			Height int    `json:"height"`
			URL    string `json:"url"`
			Width  int    `json:"width"`
		} `json:"images"`
		IsExternallyHosted   bool     `json:"is_externally_hosted"`
		IsPlayable           bool     `json:"is_playable"`
		Language             string   `json:"language"`
		Languages            []string `json:"languages"`
		Name                 string   `json:"name"`
		ReleaseDate          string   `json:"release_date"`
		ReleaseDatePrecision string   `json:"release_date_precision"`
		Show                 struct {
			AvailableMarkets []any  `json:"available_markets"`
			Copyrights       []any  `json:"copyrights"`
			Description      string `json:"description"`
			Explicit         bool   `json:"explicit"`
			ExternalUrls     struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href            string `json:"href"`
			HTMLDescription string `json:"html_description"`
			ID              string `json:"id"`
			Images          []struct {
				Height int    `json:"height"`
				URL    string `json:"url"`
				Width  int    `json:"width"`
			} `json:"images"`
			IsExternallyHosted bool     `json:"is_externally_hosted"`
			Languages          []string `json:"languages"`
			MediaType          string   `json:"media_type"`
			Name               string   `json:"name"`
			Publisher          string   `json:"publisher"`
			TotalEpisodes      int      `json:"total_episodes"`
			Type               string   `json:"type"`
			URI                string   `json:"uri"`
		} `json:"show"`
		Type string `json:"type"`
		URI  string `json:"uri"`
	} `json:"currently_playing"`
	Queue []struct {
		AudioPreviewURL string `json:"audio_preview_url"`
		Description     string `json:"description"`
		DurationMs      int    `json:"duration_ms"`
		Explicit        bool   `json:"explicit"`
		ExternalUrls    struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href            string `json:"href"`
		HTMLDescription string `json:"html_description"`
		ID              string `json:"id"`
		Images          []struct {
			Height int    `json:"height"`
			URL    string `json:"url"`
			Width  int    `json:"width"`
		} `json:"images"`
		IsExternallyHosted   bool     `json:"is_externally_hosted"`
		IsPlayable           bool     `json:"is_playable"`
		Language             string   `json:"language"`
		Languages            []string `json:"languages"`
		Name                 string   `json:"name"`
		ReleaseDate          string   `json:"release_date"`
		ReleaseDatePrecision string   `json:"release_date_precision"`
		Show                 struct {
			AvailableMarkets []any  `json:"available_markets"`
			Copyrights       []any  `json:"copyrights"`
			Description      string `json:"description"`
			Explicit         bool   `json:"explicit"`
			ExternalUrls     struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href            string `json:"href"`
			HTMLDescription string `json:"html_description"`
			ID              string `json:"id"`
			Images          []struct {
				Height int    `json:"height"`
				URL    string `json:"url"`
				Width  int    `json:"width"`
			} `json:"images"`
			IsExternallyHosted bool     `json:"is_externally_hosted"`
			Languages          []string `json:"languages"`
			MediaType          string   `json:"media_type"`
			Name               string   `json:"name"`
			Publisher          string   `json:"publisher"`
			TotalEpisodes      int      `json:"total_episodes"`
			Type               string   `json:"type"`
			URI                string   `json:"uri"`
		} `json:"show"`
		Type string `json:"type"`
		URI  string `json:"uri"`
	} `json:"queue"`
}

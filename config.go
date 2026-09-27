package main

import (
	"encoding/json"
	"os"
)

var (
	BuildClientID   = ""
	BuildDetails    = ""
	BuildState      = ""
	BuildLargeImage = ""
	BuildLargeText  = ""
	BuildSmallImage = ""
	BuildSmallText  = ""
	BuildWebUrl     = ""
	BuildDiscordUrl = ""
)

type ButtonConfig struct {
	Label string `json:"label"`
	Url   string `json:"url"`
}

type Config struct {
	ClientID              string         `json:"client_id"`
	Details               string         `json:"details"`
	State                 string         `json:"state"`
	LargeImage            string         `json:"large_image"`
	LargeText             string         `json:"large_text"`
	SmallImage            string         `json:"small_image"`
	SmallText             string         `json:"small_text"`
	Buttons               []ButtonConfig `json:"buttons"`
	GameExecutable        string         `json:"game_executable"`
	GameArguments         []string       `json:"game_arguments"`
	PatcherExecutable     string         `json:"patcher_executable"`
	PatcherArguments      []string       `json:"patcher_arguments"`
	LaunchPatcherFirst    bool           `json:"launch_patcher_first"`
	UpdateIntervalSeconds int            `json:"update_interval_seconds"`
}

func DefaultConfig() Config {
	clientID := "YOUR_DISCORD_APPLICATION_ID"
	if BuildClientID != "" {
		clientID = BuildClientID
	}

	details := "Perfect World v1.4.6"
	if BuildDetails != "" {
		details = BuildDetails
	}

	state := "Playing on Realm of Chaos"
	if BuildState != "" {
		state = BuildState
	}

	largeImage := "logo_roc"
	if BuildLargeImage != "" {
		largeImage = BuildLargeImage
	}

	largeText := "Realm of Chaos - Sirens of War"
	if BuildLargeText != "" {
		largeText = BuildLargeText
	}

	smallImage := "pwi"
	if BuildSmallImage != "" {
		smallImage = BuildSmallImage
	}

	smallText := "v1.4.6 build 2305"
	if BuildSmallText != "" {
		smallText = BuildSmallText
	}

	webUrl := "https://your-server-website.com"
	if BuildWebUrl != "" {
		webUrl = BuildWebUrl
	}

	discordUrl := "https://discord.gg/your-discord"
	if BuildDiscordUrl != "" {
		discordUrl = BuildDiscordUrl
	}

	return Config{
		ClientID:   clientID,
		Details:    details,
		State:      state,
		LargeImage: largeImage,
		LargeText:  largeText,
		SmallImage: smallImage,
		SmallText:  smallText,
		Buttons: []ButtonConfig{
			{
				Label: "🌐 Website",
				Url:   webUrl,
			},
			{
				Label: "💬 Discord Server",
				Url:   discordUrl,
			},
		},
		GameExecutable:        "element/elementclient.exe",
		GameArguments:         []string{"game:cpw", "console:1"},
		PatcherExecutable:     "patcher/patcher.exe",
		PatcherArguments:      []string{},
		LaunchPatcherFirst:    true,
		UpdateIntervalSeconds: 15,
	}
}

func LoadConfig(filename string) Config {
	cfg := DefaultConfig()
	if filename != "" {
		if data, err := os.ReadFile(filename); err == nil {
			_ = json.Unmarshal(data, &cfg)
		}
	}
	return cfg
}

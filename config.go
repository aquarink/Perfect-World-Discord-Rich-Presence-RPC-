package main

import (
	"encoding/json"
	"os"
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
	return Config{
		ClientID:   "YOUR_DISCORD_APPLICATION_ID",
		Details:    "Perfect World v1.4.6",
		State:      "Playing on Realm of Chaos",
		LargeImage: "logo_roc",
		LargeText:  "Realm of Chaos - Sirens of War",
		SmallImage: "pwi",
		SmallText:  "v1.4.6 build 2305",
		Buttons: []ButtonConfig{
			{
				Label: "🌐 Website",
				Url:   "https://your-server-website.com",
			},
			{
				Label: "💬 Discord Server",
				Url:   "https://discord.gg/your-discord",
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
	data, err := os.ReadFile(filename)
	if err == nil {
		_ = json.Unmarshal(data, &cfg)
	}
	return cfg
}

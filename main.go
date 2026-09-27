package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func findExecutable(configuredPath string, defaults []string) (string, error) {
	var candidates []string
	if configuredPath != "" {
		candidates = append(candidates, configuredPath)
	}
	candidates = append(candidates, defaults...)

	var fullCandidates []string
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		for _, c := range candidates {
			fullCandidates = append(fullCandidates, filepath.Join(exeDir, c))
		}
	}
	fullCandidates = append(fullCandidates, candidates...)

	for _, p := range fullCandidates {
		if p == "" {
			continue
		}
		absPath, err := filepath.Abs(p)
		if err == nil {
			if info, err := os.Stat(absPath); err == nil && !info.IsDir() {
				return absPath, nil
			}
		}
	}
	return "", fmt.Errorf("file executable tidak ditemukan")
}

func findConfigFile() string {
	candidates := []string{
		"config/config.json",
		"config/discord.json",
		"launcher/config.json",
		"config.json",
	}

	var fullCandidates []string
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		for _, c := range candidates {
			fullCandidates = append(fullCandidates, filepath.Join(exeDir, c))
		}
	}
	fullCandidates = append(fullCandidates, candidates...)

	for _, p := range fullCandidates {
		absPath, err := filepath.Abs(p)
		if err == nil {
			if info, err := os.Stat(absPath); err == nil && !info.IsDir() {
				return absPath
			}
		}
	}
	return ""
}

func main() {
	// 1. Single Instance Protection (prevents duplicate launchers / double clients)
	mutexHandle, isFirst := acquireSingleInstanceMutex("RealmOfChaos_Launcher_Mutex")
	if !isFirst {
		// Another instance is already active. Exit silently.
		return
	}
	defer releaseSingleInstanceMutex(mutexHandle)

	// 2. Load Configuration (checks config/config.json, config/discord.json, launcher/, or root, or embedded)
	cfgPath := findConfigFile()
	cfg := LoadConfig(cfgPath)

	// 3. Locate Executables (Patcher & Game Client)
	patcherPath, patcherErr := findExecutable(cfg.PatcherExecutable, []string{"patcher/patcher.exe", "patcher.exe"})
	gamePath, gameErr := findExecutable(cfg.GameExecutable, []string{"element/elementclient.exe", "elementclient.exe"})

	var targetLaunchPath string
	var targetLaunchDir string
	var targetLaunchArgs []string
	isLaunchingPatcher := false

	if cfg.LaunchPatcherFirst && patcherErr == nil {
		targetLaunchPath = patcherPath
		targetLaunchDir = filepath.Dir(patcherPath)
		targetLaunchArgs = cfg.PatcherArguments
		isLaunchingPatcher = true
	} else if gameErr == nil {
		targetLaunchPath = gamePath
		targetLaunchDir = filepath.Dir(gamePath)
		targetLaunchArgs = cfg.GameArguments
		isLaunchingPatcher = false
	} else {
		showAlert(
			"Realm of Chaos - Error",
			"Gagal memulai game!\n\nFile patcher ('patcher/patcher.exe') atau game client ('element/elementclient.exe') tidak ditemukan.\n"+
				"Pastikan file launcher ini diletakkan di folder utama game Perfect World Anda (sejajar dengan folder 'element' dan 'patcher').",
		)
		os.Exit(1)
	}

	// 4. Launch Target Process (supports automatic UAC elevation)
	if err := launchProcess(targetLaunchPath, targetLaunchDir, targetLaunchArgs); err != nil {
		showAlert(
			"Realm of Chaos - Error",
			fmt.Sprintf("Gagal menjalankan %s:\n%v", filepath.Base(targetLaunchPath), err),
		)
		os.Exit(1)
	}

	launchTime := time.Now()

	// On non-windows platforms, exit early (stubs)
	if runtime.GOOS != "windows" {
		return
	}

	// 5. Prepare Discord Activity Template
	var buttons []ActivityButton
	for _, b := range cfg.Buttons {
		if b.Label != "" && b.Url != "" {
			buttons = append(buttons, ActivityButton{
				Label: b.Label,
				Url:   b.Url,
			})
		}
	}

	discord := NewDiscordClient(cfg.ClientID)
	discordConnected := false
	gameHasRun := false
	patcherExitedTime := time.Time{}
	lastRpcUpdate := time.Time{}

	updateInterval := time.Duration(cfg.UpdateIntervalSeconds) * time.Second
	if updateInterval < 5*time.Second {
		updateInterval = 15 * time.Second
	}

	var activity Activity

	// 6. Process Monitoring Loop
	// Discord Rich Presence ONLY activates when elementclient.exe is running!
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		gameCount := countRunningProcesses("elementclient.exe")

		if gameCount > 0 {
			// Elementclient is running
			if !gameHasRun {
				gameHasRun = true
				startTime := time.Now().Unix()

				activity = Activity{
					Details: cfg.Details,
					State:   cfg.State,
					Timestamps: ActivityTimestamps{
						Start: startTime,
					},
					Assets: ActivityAssets{
						LargeImage: cfg.LargeImage,
						LargeText:  cfg.LargeText,
						SmallImage: cfg.SmallImage,
						SmallText:  cfg.SmallText,
					},
					Buttons: buttons,
				}
			}

			if !discordConnected {
				if err := discord.Connect(); err == nil {
					discordConnected = true
					_ = discord.SetActivity(activity)
					lastRpcUpdate = time.Now()
				}
			} else {
				if time.Since(lastRpcUpdate) >= updateInterval {
					_ = discord.SetActivity(activity)
					lastRpcUpdate = time.Now()
				}
			}
		} else {
			// gameCount == 0 (elementclient.exe is not currently running)
			if gameHasRun {
				// The game was running, but player has now exited all game clients!
				if discordConnected {
					discord.Close()
					discordConnected = false
				}
				break
			}

			// Game has not started yet
			if isLaunchingPatcher {
				patcherCount := countRunningProcesses(filepath.Base(patcherPath))
				if patcherCount > 0 {
					// Patcher is still active; continue waiting
					patcherExitedTime = time.Time{}
				} else {
					// Patcher is closed
					if patcherExitedTime.IsZero() {
						patcherExitedTime = time.Now()
					}
					// Allow up to 25 seconds grace period for elementclient.exe to initialize after patcher closes
					if time.Since(patcherExitedTime) > 25*time.Second {
						// Patcher was closed without starting the game
						break
					}
				}
			} else {
				// Direct game launch mode: allow 15 seconds for process to register
				if time.Since(launchTime) > 15*time.Second {
					break
				}
			}
		}
	}

	if discordConnected {
		discord.Close()
	}
	time.Sleep(300 * time.Millisecond)
}

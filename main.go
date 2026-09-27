package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
	"unsafe"
)

func showNativeAlert(title, message string) {
	if runtime.GOOS == "windows" {
		user32 := syscall.NewLazyDLL("user32.dll")
		messageBox := user32.NewProc("MessageBoxW")

		titlePtr, _ := syscall.UTF16PtrFromString(title)
		msgPtr, _ := syscall.UTF16PtrFromString(message)

		// MB_OK | MB_ICONEXCLAMATION = 0x00000000 | 0x00000030 = 0x30
		messageBox.Call(0, uintptr(unsafe.Pointer(msgPtr)), uintptr(unsafe.Pointer(titlePtr)), 0x30)
	} else {
		fmt.Printf("[%s] %s\n", title, message)
	}
}

func findGameExecutable(configuredPath string) (string, error) {
	candidates := []string{
		configuredPath,
		"element/elementclient.exe",
		"elementclient.exe",
	}

	// Also check relative to where the launcher executable is located
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		candidates = append(candidates,
			filepath.Join(exeDir, "element", "elementclient.exe"),
			filepath.Join(exeDir, "elementclient.exe"),
			filepath.Join(exeDir, configuredPath),
		)
	}

	for _, path := range candidates {
		if path == "" {
			continue
		}
		absPath, err := filepath.Abs(path)
		if err == nil {
			if info, err := os.Stat(absPath); err == nil && !info.IsDir() {
				return absPath, nil
			}
		}
	}

	return "", fmt.Errorf("file 'element/elementclient.exe' tidak ditemukan")
}

func main() {
	// 1. Load Configuration
	cfgPath := "config.json"
	if exePath, err := os.Executable(); err == nil {
		localCfg := filepath.Join(filepath.Dir(exePath), "config.json")
		if _, err := os.Stat(localCfg); err == nil {
			cfgPath = localCfg
		}
	}
	cfg := LoadConfig(cfgPath)

	// 2. Find Game Executable
	gamePath, err := findGameExecutable(cfg.GameExecutable)
	if err != nil {
		showNativeAlert(
			"Realm of Chaos - Error",
			"Gagal memulai game!\n\nFile 'element/elementclient.exe' tidak ditemukan.\n"+
				"Pastikan file launcher ini diletakkan di folder utama game Perfect World Anda (sejajar dengan folder 'element').",
		)
		os.Exit(1)
	}

	// 3. Launch Game Process with normal visible window
	gameDir := filepath.Dir(gamePath)
	cmd := exec.Command(gamePath, cfg.GameArguments...)
	cmd.Dir = gameDir

	// Note: DO NOT set HideWindow: true because elementclient.exe is a DirectX GUI game!
	// Hiding window prevents the game window from rendering.

	if err := cmd.Start(); err != nil {
		showNativeAlert(
			"Realm of Chaos - Error",
			fmt.Sprintf("Gagal menjalankan game client:\n%v", err),
		)
		os.Exit(1)
	}

	// 4. Initialize Discord RPC in background
	discord := NewDiscordClient(cfg.ClientID)
	startTime := time.Now().Unix()

	var buttons []ActivityButton
	for _, b := range cfg.Buttons {
		if b.Label != "" && b.Url != "" {
			buttons = append(buttons, ActivityButton{
				Label: b.Label,
				Url:   b.Url,
			})
		}
	}

	activity := Activity{
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

	stopRPC := make(chan struct{})

	// Goroutine to maintain connection and keepalive
	go func() {
		connected := false
		for {
			select {
			case <-stopRPC:
				if connected {
					discord.Close()
				}
				return
			default:
				if !connected {
					if err := discord.Connect(); err == nil {
						connected = true
						_ = discord.SetActivity(activity)
					}
				} else {
					_ = discord.SetActivity(activity)
				}

				interval := cfg.UpdateIntervalSeconds
				if interval < 5 {
					interval = 15
				}
				time.Sleep(time.Duration(interval) * time.Second)
			}
		}
	}()

	// 5. Wait for the game to exit
	_ = cmd.Wait()

	// 6. Clean up
	close(stopRPC)
	discord.Close()
	time.Sleep(500 * time.Millisecond)
}

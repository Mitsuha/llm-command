package main

/*
#include <unistd.h>
*/
import "C"

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type SessionEntry struct {
	Timestamp    time.Time `json:"timestamp"`
	UserQuery    string    `json:"user_query"`
	AICommand    string    `json:"ai_command"`
	UserAccepted bool      `json:"user_accepted"`
}

func getSessionDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "/tmp/.llm-command/plz"
	}
	return filepath.Join(homeDir, ".llm-command", "plz")
}

func hashString(s string) string {
	b := []byte(s)
	if len(b) > 32 {
		b = b[:32]
	}
	return hex.EncodeToString(b)
}

// 获取终端标识符（跨平台）
func getTerminalID() string {
	switch runtime.GOOS {
	case "linux", "darwin":
		// Linux/macOS: 调用 C 函数获取 tty name
		fd := C.int(os.Stdin.Fd())

		if name := C.ttyname(fd); name != nil {
			return C.GoString(name)
		}

	case "windows":
		// Windows: 使用 WT_SESSION (Windows Terminal) 或其他标识
		if wtSession := os.Getenv("WT_SESSION"); wtSession != "" {
			return wtSession
		}
		// CMD/PowerShell: 使用窗口标题哈希
		if title := os.Getenv("PROMPT"); title != "" {
			return title
		}
	}

	return ""
}

func getSessionFile() string {
	sessionDir := getSessionDir()
	terminalID := hashString(getTerminalID())
	return filepath.Join(sessionDir, terminalID+".json")
}

func loadSessionHistory() []SessionEntry {
	sessionFile := getSessionFile()
	data, err := os.ReadFile(sessionFile)
	if err != nil {
		return []SessionEntry{}
	}

	var history []SessionEntry
	if err := json.Unmarshal(data, &history); err != nil {
		return []SessionEntry{}
	}
	return history
}

func saveSessionEntry(userQuery, aiCommand string, accepted bool) {
	sessionDir := getSessionDir()
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return // Silently fail if can't create directory
	}

	history := loadSessionHistory()
	entry := SessionEntry{
		Timestamp:    time.Now(),
		UserQuery:    userQuery,
		AICommand:    aiCommand,
		UserAccepted: accepted,
	}
	history = append(history, entry)

	// Keep only last 50 entries per session
	if len(history) > 50 {
		history = history[len(history)-50:]
	}

	sessionFile := getSessionFile()
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(sessionFile, data, 0644)
}

func cleanupOldSessions() {
	sessionDir := getSessionDir()
	files, err := os.ReadDir(sessionDir)
	if err != nil {
		return
	}

	cutoff := time.Now().AddDate(0, 0, -3) // 3 days ago
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			filePath := filepath.Join(sessionDir, file.Name())
			if info, err := file.Info(); err == nil {
				if info.ModTime().Before(cutoff) {
					_ = os.Remove(filePath)
				}
			}
		}
	}
}

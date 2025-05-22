package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Settings struct {
	RemoveComments bool   `json:"remove_comments"`
	IncludeTests   bool   `json:"include_tests"`
	MinifyOutput   bool   `json:"minify_output"`
	LastSrcPath    string `json:"last_src_path"`
	LastDestPath   string `json:"last_dest_path"`
}

func LoadSettings() *Settings {
	settings := &Settings{
		RemoveComments: true,
		IncludeTests:   false,
		MinifyOutput:   true,
	}

	configPath := getConfigPath()
	if data, err := os.ReadFile(configPath); err == nil {
		json.Unmarshal(data, settings)
	}

	return settings
}

func (s *Settings) Save() error {
	configPath := getConfigPath()
	os.MkdirAll(filepath.Dir(configPath), 0755)

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func getConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "go-context-generator", "settings.json")
}

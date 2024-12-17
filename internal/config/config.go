package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "", fmt.Errorf("Could not find $HOME: %w", err)
    }
    path := filepath.Join(homeDir, "bootdev", "gator", configFileName)
    return path, nil
}

func write(cfg Config) error {
    filePath, err := getConfigFilePath()
    if err != nil {
        return fmt.Errorf("Could not write to file: %w", err)
    }

    f, err := os.OpenFile(filePath, os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("Could not open file: %w", err)
    }
    defer f.Close()

    data, err := json.Marshal(cfg)
    if err != nil {
        return fmt.Errorf("Could not Marshal config: %w", err)
    }

    _, err = f.Write(data)
    if err != nil {
        return fmt.Errorf("Something went wrong when saving file: %w", err)
    }

    return nil
}

func Read() (Config, error) {
    filePath, err := getConfigFilePath()
    if err != nil {
        return Config{}, err
    }

    data, err := os.ReadFile(filePath)
    if err != nil {
        return Config{}, err
    }

    var config Config
    err = json.Unmarshal(data, &config)
    if err != nil {
        return Config{}, fmt.Errorf("Could not read config file: %w", err)
    }

    return config, nil
}



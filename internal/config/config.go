package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)



type Config struct {
    DBUrl		string `json:"db_url"`
    CurrentUserName	string `json:"current_user_name"`
}



func Read() (Config, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
	return Config{}, err
    }

    fullPath := filepath.Join(homeDir, ".gatorconfig.json")


    data, err := os.ReadFile(fullPath)
    if err != nil {
	return Config{}, err
    }

    var config Config

    err = json.Unmarshal(data, &config)
    if err != nil {
	return Config{}, err
    }

    return config, nil
}


func (cfg *Config) SetUser(username string) error {
    cfg.CurrentUserName = username

    data, err := json.Marshal(cfg)
    if err != nil {
	return err
    }

    homeDir, err := os.UserHomeDir()
    if err != nil {
	return err
    }

    fullPath := filepath.Join(homeDir, ".gatorconfig.json")

    err = os.WriteFile(fullPath, data, 0644)
    if err != nil {
	return err
    }

    return nil
}
























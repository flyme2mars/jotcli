package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func InitConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// 1. Set Defaults
	defaultDBPath := filepath.Join(home, ".jot.db")
	viper.SetDefault("database", defaultDBPath)
	viper.SetDefault("editor", "vim")

	// 2. Set config file details
	viper.SetConfigName(".jotcli") // Name: ~/.jotcli.yaml
	viper.SetConfigType("yaml")
	viper.AddConfigPath(home)

	// 3. Enable environment variables (e.g., JOT_DATABASE)
	viper.SetEnvPrefix("JOT")
	viper.AutomaticEnv()

	// 4. Read config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
		// It's okay if the config file is missing
	}

	return nil
}

func GetDBPath() string {
	return viper.GetString("database")
}

func GetEditor() string {
	// Priority: 1. Environment Variable ($EDITOR) 2. Config File (.jotcli.yaml) 3. Default (vim)
	envEditor := os.Getenv("EDITOR")
	if envEditor != "" {
		return envEditor
	}
	return viper.GetString("editor")
}

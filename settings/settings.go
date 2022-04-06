package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Settings struct {
	DatabaseUrl   string
	SigningSecret string
	Port          string
}

func (settings *Settings) Serialize() []byte {
	b, err := json.Marshal(&settings)
	if err != nil {
		fmt.Println(err)
	}
	return b
}

func (settings *Settings) Deserialize(bytes []byte) {
	err := json.Unmarshal(bytes, &settings)
	if err != nil {
		fmt.Println(err)
	}
}

func (settings *Settings) IsValid() bool {
	return settings.SigningSecret != ""
}

func LoadSettings() *Settings {
	absolutePath, err := filepath.Abs("./settings/settings.json")
	if err != nil {
		fmt.Println(err)
	}
	if _, err := os.Stat(absolutePath); errors.Is(err, os.ErrNotExist) {
		return loadFromEnvironmentVariables()
	}
	bytes, err := os.ReadFile(absolutePath)
	if err != nil {
		fmt.Println(err)
	}
	settings := Settings{}
	settings.Deserialize(bytes)
	return &settings
}

func loadFromEnvironmentVariables() *Settings {
	settings := Settings{
		DatabaseUrl:   os.Getenv("DATABASE_URL"),
		SigningSecret: os.Getenv("SIGNING_SECRET"),
		Port:          os.Getenv("PORT"),
	}
	return &settings
}

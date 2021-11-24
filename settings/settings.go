package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Settings struct {
	Driver          string
	Host            string
	User            string
	Password        string
	AllowedLogin    string
	AllowedPassword string
	SigningSecret   string
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
	return settings.Driver != "" && settings.Host != "" && settings.User != "" && settings.Password != ""
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
		Driver:          os.Getenv("Spendon_driver"),
		Host:            os.Getenv("Spendon_host"),
		User:            os.Getenv("Spendon_user"),
		Password:        os.Getenv("Spendon_password"),
		AllowedLogin:    os.Getenv("ALLOWED_LOGIN"),
		AllowedPassword: os.Getenv("ALLOWED_PASSWORD"),
		SigningSecret:   os.Getenv("SIGNING_SECRET"),
	}
	return &settings
}

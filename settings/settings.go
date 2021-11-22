package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Settings struct {
	Driver   string
	Host     string
	User     string
	Password string
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
	bytes, err := os.ReadFile(absolutePath)
	if err != nil {
		fmt.Println(err)
	}
	settings := Settings{}
	settings.Deserialize(bytes)
	return &settings
}

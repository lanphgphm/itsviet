package config 

import (
	"encoding/json"
	"os"

	"itsvietbaby/internal/processor"
)

type Config struct {
	InputMethod			processor.InputMethod  	`json:"input_method"`
	BrowserClassNames	[]string				`json:"browser_class_names"`
}

func Load(path string) (*Config, error) {
	config := &Config {
		InputMethod: processor.Telex, 
		BrowserClassNames: []string {
			"Firefox", 
			"Chromium",
			"Chrome", 
			"chromium-browser", 
			"firefox",
		},
	}

	if path == "" {
		return config, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}
package config

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"gopkg.in/yaml.v3"
)

// Config represents our overall dotfiles config.
type Config struct {
	Environments []string             `yaml:"environments"`
	Globals      map[string]any       `yaml:"globals"`
	AppConfigs   map[string]AppConfig `yaml:"apps"`
}

// AppConfig represents the configuration for a particular app.
type AppConfig struct {
	Enabled bool           `yaml:"enabled"`
	Vars    map[string]any `yaml:"vars"`
}

// ConfigTemplateData represents the template data that is supported in our config.
type ConfigTemplateData struct {
	Environment string
}

// NewConfig creates a new configuration from a YAML file.
func NewConfig(configPath, environmentType string) (*Config, error) {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	// Parse and render our config file
	tmpl, err := template.New("config").Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to load config as template: %v", err)
	}
	var renderedContent bytes.Buffer
	err = tmpl.Execute(&renderedContent, ConfigTemplateData{Environment: environmentType})
	if err != nil {
		return nil, fmt.Errorf("failed to render template: %v", err)
	}

	var c Config
	if err := yaml.Unmarshal(renderedContent.Bytes(), &c); err != nil {
		return nil, fmt.Errorf("parsing config YAML: %w", err)
	}
	return &c, nil
}

// AppConfig returns the AppConfig for a given application.
// If no config was found we will instead return a default configuration.
func (c *Config) AppConfig(name string) AppConfig {
	if appConfig, ok := c.AppConfigs[name]; ok {
		return appConfig
	}
	fmt.Printf("NOTE: Using default AppConfig for %q\n", name)
	return AppConfig{Enabled: true}
}

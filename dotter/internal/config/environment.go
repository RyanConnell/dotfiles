package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// EnvironmentConfig caches the type of environment we are installed in to avoid prompting the user
// on future installations.
type EnvironmentConfig struct {
	EnvironmentType string `yaml:"environmentType"`
}

// DetermineEnvironmentType returns the type of environment we are running in.
func DetermineEnvironmentType(configPath string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	envConfigPath := filepath.Join(homeDir, ".dotfiles.yaml")

	// If we've already cached the environmentType then return that.
	ec, err := readEnvironmentConfig(envConfigPath)
	if err == nil && ec != nil && ec.EnvironmentType != "" {
		fmt.Printf("Environment: %q (Persisted to %q)\n", ec.EnvironmentType, envConfigPath)
		return ec.EnvironmentType, nil
	}

	// Otherwise, figure out which environments are supported and prompt the user to select one.
	environmentType, err := promptEnvironmentType(configPath)
	if err != nil {
		return "", err
	}

	// Persist this choice to file so we can read it next time this is run.
	ec = &EnvironmentConfig{EnvironmentType: environmentType}
	yamlContent, err := yaml.Marshal(ec)
	if err != nil {
		return "", fmt.Errorf("failed to marshal new config: %w", err)
	}
	fmt.Printf("Writing choice to %q\n", envConfigPath)
	err = os.WriteFile(envConfigPath, yamlContent, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write %s: %w", envConfigPath, err)
	}

	return environmentType, nil
}

// promptEnvironmentType prompts the user to select one of the supported environments.
func promptEnvironmentType(configPath string) (string, error) {
	// Note that the environmentType is empty when loading the config here since we don't know
	// what it is yet. The config will be reloaded later with the correct configPath.
	config, err := NewConfig(configPath, "")
	if err != nil {
		return "", fmt.Errorf("failed to load config for environment selection: %w", err)
	}
	if len(config.Environments) == 0 {
		return "", fmt.Errorf("unable to determine list of valid environments...")
	}

	userPrompt := "\nPlease select an environment:"
	for i, env := range config.Environments {
		userPrompt += fmt.Sprintf("\n[%d] %s", i+1, env)
	}
	userPrompt += "\n\n> "

	fmt.Print(userPrompt)
	var choice int
	_, err = fmt.Scanln(&choice)
	if err != nil || choice < 1 || choice > len(config.Environments) {
		return "", fmt.Errorf("invalid environment selection: %v", err)
	}

	return config.Environments[choice-1], nil
}

// readEnvironmentConfig reads the environmentType from our config file.
func readEnvironmentConfig(path string) (*EnvironmentConfig, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}
	var envConfig *EnvironmentConfig
	if err := yaml.Unmarshal(content, &envConfig); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", path, err)
	}
	return envConfig, nil
}

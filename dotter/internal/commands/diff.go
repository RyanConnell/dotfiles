package commands

import (
	"fmt"
	"path/filepath"

	"github.com/RyanConnell/dotfiles/dotter/internal/apps"
	"github.com/RyanConnell/dotfiles/dotter/internal/config"
)

type Differ struct {
	Config string
	Apps   string
	Output string
}

func NewDiffer(config, apps, output string) *Differ {
	return &Differ{Config: config, Apps: apps, Output: output}
}

func (cmd *Differ) Run() error {
	envType, err := config.DetermineEnvironmentType(cmd.Config)
	if err != nil {
		return fmt.Errorf("failed to determine environment type: %w", err)
	}

	cfg, err := config.NewConfig(cmd.Config, envType)
	if err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	applications, err := apps.DiscoverApps(cmd.Apps, envType, cfg)
	if err != nil {
		return fmt.Errorf("failed to load apps from directory: %w", err)
	}

	diff, err := fullAppDiff(applications, cmd.Apps, cmd.Output)
	if err != nil {
		return err
	}

	if diff != "" {
		fmt.Println(diff)
	} else {
		fmt.Println("No differences found")
	}

	return nil
}

func fullAppDiff(applications []*apps.App, appDir, outputDir string) (string, error) {
	var combinedDiff string
	for _, app := range applications {
		source := filepath.Join(appDir, app.Name)
		target := filepath.Join(outputDir, app.Name)
		diff, err := app.Differences(source, target)
		if err != nil {
			return "", fmt.Errorf("failed to capture differences between %q and %q: %v",
				source, target, err)
		}
		if diff == "" {
			continue
		}
	}
	return combinedDiff, nil
}

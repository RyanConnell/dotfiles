package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/RyanConnell/dotfiles/dotter/internal/apps"
	"github.com/RyanConnell/dotfiles/dotter/internal/config"
)

type Installer struct {
	Config string
	Apps   string
	Output string
}

func NewInstaller(config, apps, output string) *Installer {
	return &Installer{Config: config, Apps: apps, Output: output}
}

func (cmd *Installer) Run() error {
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

	// Do a quick pre-check at the beginning of the run to make sure we aren't overwriting any
	// newer content from the build directory.
	diff, err := fullAppDiff(applications, cmd.Apps, cmd.Output)
	if diff != "" {
		fmt.Printf("Found differences between %q and %q:\n\n%s\n",
			cmd.Apps, cmd.Output, diff)
		if !userWantsToContinue() {
			return fmt.Errorf("User aborted operation")
		}
	}

	failures := make(map[string]error)
	fail := func(appName, message string, args ...any) {
		fmt.Printf("[%s]: ERROR: %s; Skipping\n", appName, fmt.Sprintf(message, args...))
		failures[appName] = fmt.Errorf("%s", fmt.Sprintf(message, args...))
	}

	for _, app := range applications {
		fmt.Println("\n---------------------------------")
		fmt.Printf("[%s]: Running pre.sh...\n", app.Name)
		if err := app.MaybeRunScript("pre.sh"); err != nil {
			fail(app.Name, "failed running pre.sh: %v", err)
			continue
		}

		fmt.Printf("[%s]: Stowing package...\n", app.Name)
		if err := app.Render(cmd.Output); err != nil {
			fail(app.Name, "failed rendering package: %v", err)
			continue
		}
		// Check if there are any conflicts that need interactive handling
		conflicts, err := app.Stow(cmd.Output, false)
		if len(conflicts) != 0 {
			err = cmd.handleConflicts(app, cmd.Output, conflicts)
			if err != nil {
				fail(app.Name, "failed conflict resolution: %v", err)
				continue
			}
		} else if err != nil {
			fail(app.Name, "initial stow failed: %v", err)
			continue
		}

		fmt.Printf("[%s]: Running post.sh...\n", app.Name)
		if err := app.MaybeRunScript("post.sh"); err != nil {
			fail(app.Name, "failed running post.sh: %v", err)
			continue
		}
	}

	fmt.Printf("\n=================================\n%d Failures occurred\n", len(failures))
	for appName, failure := range failures {
		fmt.Printf("- %q: %v\n", appName, failure)
	}

	return nil
}

// handleConflicts performs an interactive resolution of stowing conflicts by presenting
// a diff and asking the user if they wish to adopt the existing files.
func (cmd *Installer) handleConflicts(app *apps.App, outputDir string, conflicts []string) error {
	fmt.Printf("[%s]: Found %d conflicts for package: %v\n", app.Name, len(conflicts), conflicts)

	hasDifferences := false
	for _, relPath := range conflicts {
		source := filepath.Join(cmd.Output, app.Name, relPath)
		target := filepath.Join(os.Getenv("HOME"), relPath)
		diff, err := app.DiffFiles(source, target)
		if err != nil {
			return fmt.Errorf("failed to check diff for %s: %w", relPath, err)
		}
		if diff != "" {
			hasDifferences = true
			fmt.Printf("\n--- DIFF FOR: %s ---\n%s\n", relPath, diff)
		}
	}

	if !hasDifferences {
		fmt.Printf("[%s]: No differences found in conflicting files. Deleting and restowing...\n",
			app.Name)
		for _, file := range conflicts {
			err := os.Remove(filepath.Join(os.Getenv("HOME"), file))
			if err != nil {
				return err
			}
		}
		_, err := app.Stow(outputDir, true)
		return err
	}

	// Prompt the user for adoption
	if userWantsToContinue() {
		fmt.Printf("[%s]: Adopting and overwriting content...\n", app.Name)
		_, err := app.Stow(outputDir, true)
		if err != nil {
			return fmt.Errorf("adoption failed: %w", err)
		}
		return app.Render(outputDir)
	}
	return errors.New("user declined override")
}

func userWantsToContinue() bool {
	fmt.Print("Do you want to continue? (y/N): ")
	var response string
	_, _ = fmt.Scanln(&response)
	return strings.ToLower(response) == "y"
}

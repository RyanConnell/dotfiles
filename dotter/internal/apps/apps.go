package apps

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/RyanConnell/dotfiles/dotter/internal/config"
)

var ignoredFiles = map[string]struct{}{
	"pre.sh":  {},
	"post.sh": {},
}

// AppTemplateData contains the fields passed into all of our app templates.
type AppTemplateData struct {
	Environment string
	Globals     map[string]any
	Vars        map[string]any
}

// App describes an application we want to install configuration for.
type App struct {
	Name            string
	SourcePath      string
	AppTemplateData AppTemplateData
	Files           []string
}

// NewApp loads information about an application
func NewApp(name, sourcePath string, appTemplateData AppTemplateData) (*App, error) {
	files, err := findFiles(sourcePath)
	if err != nil {
		return nil, err
	}
	return &App{
		Name:            name,
		SourcePath:      sourcePath,
		AppTemplateData: appTemplateData,
		Files:           files,
	}, nil
}

// RunScript executes a script located at '/{sourcePath}/{scriptName}'
// and prefixes its output with the application name.
func (a *App) RunScript(scriptName string) error {
	scriptPath := filepath.Join(a.SourcePath, scriptName)
	cmdStr := fmt.Sprintf("%s | sed 's/^/[%s] /'", scriptPath, a.Name)
	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// MaybeRunScript attempts to run the given script but avoids failing if the
// file does not exist.
func (a *App) MaybeRunScript(scriptName string) error {
	_, err := os.Stat(filepath.Join(a.SourcePath, scriptName))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return a.RunScript(scriptName)
}

// Render copies the applications files to the output directory.
func (a *App) Render(outputDir string) error {
	var err error
	for _, file := range a.StowableFiles() {
		sourceFilePath := filepath.Join(a.SourcePath, file)
		targetFilePath := filepath.Join(outputDir, a.Name, file)
		if strings.HasSuffix(file, ".tmpl") {
			targetFilePath = strings.TrimSuffix(targetFilePath, ".tmpl")
			err = a.renderTemplate(sourceFilePath, targetFilePath)
		} else {
			err = copyFile(sourceFilePath, targetFilePath)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) renderTemplate(sourceFilePath, targetFilePath string) error {
	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	content, err := io.ReadAll(sourceFile)
	if err != nil {
		return err
	}
	tmpl, err := template.New(filepath.Base(sourceFilePath)).Parse(string(content))
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(targetFilePath), 0755); err != nil {
		return err
	}

	// Render our template
	// Write results to a file.
	sourceFileStat, err := sourceFile.Stat()
	if err != nil {
		return err
	}
	targetFile, err := os.OpenFile(targetFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, sourceFileStat.Mode())
	if err != nil {
		return err
	}

	return tmpl.Execute(targetFile, a.AppTemplateData)
}

func (a *App) StowableFiles() []string {
	var stowable []string
	for _, file := range a.Files {
		if _, ok := ignoredFiles[file]; ok {
			continue
		}
		stowable = append(stowable, file)
	}
	return stowable
}

// Stow attempts to use 'stow' to install an applications config files
func (a *App) Stow(outputDir string, adopt bool) ([]string, error) {
	if len(a.StowableFiles()) == 0 {
		// Skip running stow if we have nothing to stow.
		return nil, nil
	}

	stowArgs := []string{
		"-d", outputDir,
		"-t", os.Getenv("HOME"),
		"--ignore", "(pre|post).sh",
	}
	if adopt {
		stowArgs = append(stowArgs, "--adopt")
	}
	stowArgs = append(stowArgs, a.Name)

	cmd := exec.Command("stow", stowArgs...)
	fmt.Println("stow", stowArgs)
	var stderrBuf bytes.Buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderrBuf

	err := cmd.Run()

	// Figure out which files (if any) failed to stow.
	var conflicts []string
	re := regexp.MustCompile(`(?m)^\s+\*\s+existing target is not owned by stow:\s+(.*)$`)
	stderrStr := stderrBuf.String()
	matches := re.FindAllStringSubmatch(stderrStr, -1)
	for _, match := range matches {
		if len(match) > 1 {
			conflicts = append(conflicts, match[1])
		}
	}

	// Write the error message back into Stderr
	if err != nil {
		return conflicts, fmt.Errorf("%v: %v", err, stderrStr)
	}
	return conflicts, nil
}

// Differences returns a diff between the two directories, excluding ignored files.
func (a *App) Differences(sourceDir, targetDir string) (string, error) {
	var combinedDiff string
	for _, file := range a.StowableFiles() {
		diff, err := a.DiffFiles(filepath.Join(sourceDir, file), filepath.Join(targetDir, file))
		if err != nil {
			return "", err
		}
		if diff == "" {
			continue
		}
		combinedDiff += fmt.Sprintf("%s\n", diff)
	}
	return combinedDiff, nil
}

// DiffFiles returns the diff between two files.
func (a *App) DiffFiles(source, target string) (string, error) {
	// If we're dealing with a template we'll need to render it to a temporary file to check
	// the differences. (Since i'm lazy and don't want to do it in memory)
	if strings.HasSuffix(source, ".tmpl") {
		// Create a temporary file.
		pattern := filepath.Base(source) + ".tmp"
		tmpFile, err := os.CreateTemp(filepath.Dir(target), pattern)
		if err != nil {
			return "", fmt.Errorf("failed to create temporary file: %v", err)
		}

		// Render template to our temporary file.
		if err := a.renderTemplate(source, tmpFile.Name()); err != nil {
			return "", fmt.Errorf("failed to render template: %v", err)
		}
		defer func() {
			// Ensure we clean up the file when this function exits.
			if err := os.Remove(tmpFile.Name()); err != nil {
				fmt.Printf("WARNING: Failed to clean up temporary file %q: %v", tmpFile.Name(), err)
			}
		}()

		// Update our target and source files to match the expected ones.
		target = strings.TrimSuffix(target, ".tmpl")
		source = tmpFile.Name()
	}

	// If the target file doesn't exist at all, report it as entirely new.
	if _, err := os.Stat(target); errors.Is(err, os.ErrNotExist) {
		return fmt.Sprintf("\033[1m--- %s\n+++ %s\033[0m\n\033[1;32m<entirely new>\033[0m\n",
			target, source), nil
	}

	cmd := exec.Command("diff", "--color=always", "-u", target, source)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			return out.String(), nil
		}
		return "", err
	}

	return "", nil
}

// DiscoverApps walks the app folder to gather information about all available apps.
func DiscoverApps(path, envType string, cfg *config.Config) ([]*App, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("reading apps directory: %w", err)
	}

	var apps []*App
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		appTemplateData := AppTemplateData{
			Environment: envType,
			Globals:     cfg.Globals,
			Vars:        cfg.AppConfig(entry.Name()).Vars,
		}
		app, err := NewApp(entry.Name(), filepath.Join(path, entry.Name()), appTemplateData)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}
	return apps, nil
}

// findFiles walks a directory to find all files within it.
func findFiles(sourcePath string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(sourcePath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return err
		}
		relativePath, err := filepath.Rel(sourcePath, path)
		if err != nil {
			return err
		}
		files = append(files, relativePath)
		return err
	})
	return files, err
}

func copyFile(src, dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}

	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	info, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	destFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

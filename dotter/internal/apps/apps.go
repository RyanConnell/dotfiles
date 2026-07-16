package apps

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

var ignoredFiles = map[string]struct{}{
	"pre.sh":  {},
	"post.sh": {},
}

type App struct {
	Name       string
	SourcePath string
	Files      []string
}

// NewApp loads information about an application
func NewApp(name, sourcePath string) (*App, error) {
	files, err := findFiles(sourcePath)
	if err != nil {
		return nil, err
	}
	return &App{Name: name, SourcePath: sourcePath, Files: files}, nil
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
	for _, file := range a.Files {
		if _, ok := ignoredFiles[file]; ok {
			continue
		}
		err := copyFile(filepath.Join(a.SourcePath, file), filepath.Join(outputDir, a.Name, file))
		if err != nil {
			return err
		}
	}
	return nil
}

// Stow attempts to use 'stow' to install an applications config files
func (a *App) Stow(outputDir string, adopt bool) ([]string, error) {
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
	for _, file := range a.Files {
		if _, ok := ignoredFiles[file]; ok {
			continue
		}
		diff, err := DiffFiles(filepath.Join(sourceDir, file), filepath.Join(targetDir, file))
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

// DiscoverApps walks the app folder to gather information about all available apps.
func DiscoverApps(path string) ([]*App, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("reading apps directory: %w", err)
	}

	var apps []*App
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		app, err := NewApp(entry.Name(), filepath.Join(path, entry.Name()))
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

// DiffFiles returns the diff between two files.
func DiffFiles(source, target string) (string, error) {
	cmd := exec.Command("diff", "--color=always", "-u", source, target)
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

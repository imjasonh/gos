package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type scriptMetadata struct {
	Dependencies []string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: gos run <script.go> [args...]\n")
		os.Exit(1)
	}

	if os.Args[1] != "run" {
		fmt.Fprintf(os.Stderr, "Error: only 'run' command is supported\n")
		fmt.Fprintf(os.Stderr, "Usage: gos run <script.go> [args...]\n")
		os.Exit(1)
	}

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Error: script file required\n")
		fmt.Fprintf(os.Stderr, "Usage: gos run <script.go> [args...]\n")
		os.Exit(1)
	}

	scriptPath := os.Args[2]
	scriptArgs := os.Args[3:]

	// Read and parse the script
	metadata, scriptContent, err := parseScript(scriptPath)
	if err != nil {
		log.Fatalf("Failed to parse script: %v", err)
	}

	// Create temporary workspace
	tempDir, err := os.MkdirTemp("", "gos-*")
	if err != nil {
		log.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Generate go.mod
	if err := generateGoMod(tempDir, metadata); err != nil {
		log.Fatalf("Failed to generate go.mod: %v", err)
	}

	// Copy script to temp directory
	scriptName := filepath.Base(scriptPath)
	tempScriptPath := filepath.Join(tempDir, scriptName)
	if err := os.WriteFile(tempScriptPath, scriptContent, 0644); err != nil {
		log.Fatalf("Failed to write script: %v", err)
	}

	// Run go mod tidy
	if err := runGoModTidy(tempDir); err != nil {
		log.Fatalf("Failed to run go mod tidy: %v", err)
	}

	// Build and run the script
	if err := buildAndRun(tempDir, scriptName, scriptArgs); err != nil {
		log.Fatalf("Failed to build and run script: %v", err)
	}
}

func parseScript(scriptPath string) (*scriptMetadata, []byte, error) {
	file, err := os.Open(scriptPath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	var buf bytes.Buffer
	scanner := bufio.NewScanner(file)
	inMetadata := false
	var metadataLines []string
	metadata := &scriptMetadata{}

	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "// /// script" {
			inMetadata = true
			continue
		}

		if inMetadata && strings.TrimSpace(line) == "// ///" {
			inMetadata = false
			// Parse metadata
			if err := parseMetadata(metadataLines, metadata); err != nil {
				return nil, nil, err
			}
			continue
		}

		if inMetadata {
			// Remove comment prefix
			trimmed := strings.TrimPrefix(line, "//")
			trimmed = strings.TrimSpace(trimmed)
			if trimmed != "" {
				metadataLines = append(metadataLines, trimmed)
			}
		} else if !strings.HasPrefix(line, "#!") {
			// Skip shebang line
			buf.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return metadata, buf.Bytes(), nil
}

func parseMetadata(lines []string, metadata *scriptMetadata) error {
	// Join lines and look for dependencies
	content := strings.Join(lines, "\n")

	// Match dependencies = [ ... ] (with multiline support)
	depRegex := regexp.MustCompile(`(?s)dependencies\s*=\s*\[(.*?)\]`)
	matches := depRegex.FindStringSubmatch(content)

	if len(matches) > 1 {
		depString := matches[1]
		// Parse individual dependencies
		deps := strings.Split(depString, ",")
		for _, dep := range deps {
			dep = strings.TrimSpace(dep)
			dep = strings.Trim(dep, `"'`)
			if dep != "" {
				metadata.Dependencies = append(metadata.Dependencies, dep)
			}
		}
	}

	return nil
}

func generateGoMod(dir string, metadata *scriptMetadata) error {
	content := "module goscript\n\ngo 1.21\n"

	if len(metadata.Dependencies) > 0 {
		content += "\nrequire (\n"
		for _, dep := range metadata.Dependencies {
			parts := strings.Split(dep, "@")
			if len(parts) == 2 {
				content += fmt.Sprintf("\t%s %s\n", parts[0], parts[1])
			} else {
				content += fmt.Sprintf("\t%s latest\n", dep)
			}
		}
		content += ")\n"
	}

	return os.WriteFile(filepath.Join(dir, "go.mod"), []byte(content), 0644)
}

func runGoModTidy(dir string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func buildAndRun(dir, scriptName string, args []string) error {
	// Build the script
	binaryName := strings.TrimSuffix(scriptName, ".go")
	buildCmd := exec.Command("go", "build", "-o", binaryName, scriptName)
	buildCmd.Dir = dir
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	// Run the binary
	binaryPath := filepath.Join(dir, binaryName)
	runCmd := exec.Command(binaryPath, args...)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	runCmd.Stdin = os.Stdin

	return runCmd.Run()
}

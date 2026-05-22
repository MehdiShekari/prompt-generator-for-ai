package main

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//go:embed banner.txt
var banner string

// ANSI color codes
const (
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorGray  = "\033[90m"
	colorCyan  = "\033[36m"
)

// Config holds user settings loaded from config.json
type Config struct {
	BinaryExtensions []string `json:"binary_extensions"`
	IgnorePatterns   []string `json:"ignore_patterns"`
	IgnoreHidden     bool     `json:"ignore_hidden"`
	MaxFileSizeKB    int64    `json:"max_file_size_kb"`
}

// DefaultConfig returns a Config with sensible defaults
func DefaultConfig() Config {
	return Config{
		BinaryExtensions: []string{
			".png", ".jpg", ".jpeg", ".gif", ".bmp", ".ico",
			".ttf", ".woff", ".woff2", ".eot", ".otf", ".svg",
			".mp3", ".mp4", ".webp",
		},
		IgnorePatterns: []string{
			".git/*",
			"node_modules/*",
			"*.pyc",
			"__pycache__/*",
		},
		IgnoreHidden:  true,
		MaxFileSizeKB: 500,
	}
}

// loadConfig reads config from the given path. If the file doesn't exist,
// it creates one with default values and returns the defaults.
func loadConfig(path string) (Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		cfg := DefaultConfig()
		data, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			return cfg, fmt.Errorf("failed to marshal default config: %w", err)
		}
		if err := ioutil.WriteFile(path, data, 0644); err != nil {
			return cfg, fmt.Errorf("failed to write default config: %w", err)
		}
		fmt.Printf("%sCreated default config: %s%s\n", colorGray, path, colorReset)
		return cfg, nil
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("cannot read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("invalid config JSON: %w", err)
	}
	return cfg, nil
}

// shouldIgnore checks if a given file should be skipped based on config.
func shouldIgnore(relPath string, info os.FileInfo, cfg Config) bool {
	// Hidden files/directories (name starts with dot)
	if cfg.IgnoreHidden && strings.HasPrefix(info.Name(), ".") {
		return true
	}

	// Check glob ignore patterns
	for _, pattern := range cfg.IgnorePatterns {
		matched, err := filepath.Match(pattern, relPath)
		if err == nil && matched {
			return true
		}
		// Also handle "dirname/*" patterns
		if strings.HasSuffix(pattern, "/*") {
			dirPattern := strings.TrimSuffix(pattern, "/*")
			if matched, _ := filepath.Match(dirPattern, relPath); matched {
				return true
			}
		}
	}

	// Directory: don't ignore by size/extension (handled in Walk)
	if info.IsDir() {
		return false
	}

	// File size limit
	if cfg.MaxFileSizeKB > 0 && info.Size() > cfg.MaxFileSizeKB*1024 {
		return true
	}

	// Binary extension check
	ext := strings.ToLower(filepath.Ext(relPath))
	for _, be := range cfg.BinaryExtensions {
		if ext == strings.ToLower(be) {
			return true
		}
	}

	return false
}

// isBinaryContent checks if data appears to be binary (contains null byte).
func isBinaryContent(data []byte) bool {
	for _, b := range data {
		if b == 0 {
			return true
		}
	}
	return false
}

// buildFileTree recursively scans the directory and returns structure lines and file contents.
func buildFileTree(root string, cfg Config) (structure []string, contents []string, err error) {
	err = filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			fmt.Printf("%sWarning: cannot access %s: %v%s\n", colorRed, path, walkErr, colorReset)
			return nil
		}

		relPath, _ := filepath.Rel(root, path)
		relPath = filepath.ToSlash(relPath) // normalize slashes

		// Root directory
		if relPath == "." {
			structure = append(structure, filepath.Base(root)+":")
			return nil
		}

		// Determine indent
		depth := strings.Count(relPath, "/")
		prefix := strings.Repeat("--------", depth)

		// If directory
		if info.IsDir() {
			// Check if directory should be skipped entirely
			for _, pattern := range cfg.IgnorePatterns {
				matched, _ := filepath.Match(pattern, relPath)
				if matched {
					return filepath.SkipDir
				}
				if strings.HasSuffix(pattern, "/*") {
					dirPattern := strings.TrimSuffix(pattern, "/*")
					if matched, _ := filepath.Match(dirPattern, relPath); matched {
						return filepath.SkipDir
					}
				}
			}
			structure = append(structure, prefix+info.Name()+":")
			return nil
		}

		// Regular file: add to structure
		structure = append(structure, prefix+info.Name())

		// Check if file should be ignored (binary, hidden, size, pattern)
		if shouldIgnore(relPath, info, cfg) {
			fmt.Printf("%sSkipping: %s%s\n", colorRed, relPath, colorReset)
			return nil
		}

		// Read file content
		data, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Printf("%sError reading %s: %v%s\n", colorRed, relPath, err, colorReset)
			contents = append(contents, fmt.Sprintf("\n--- %s ---\n[Could not read file: %v]", relPath, err))
			return nil
		}

		// Automatic binary content detection
		if isBinaryContent(data) {
			fmt.Printf("%sSkipping binary file: %s%s\n", colorRed, relPath, colorReset)
			return nil
		}

		fmt.Printf("%sAdded: %s%s\n", colorGreen, relPath, colorReset)
		contents = append(contents, fmt.Sprintf("\n--- %s ---\n%s", relPath, string(data)))
		return nil
	})
	return
}

// resolveTargetDir takes user input and tries to find the directory,
// first relative to executable, then to current working directory.
func resolveTargetDir(input string, exeDir string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		// Default: current working directory
		return os.Getwd()
	}

	// Absolute path: validate and use
	if filepath.IsAbs(input) {
		info, err := os.Stat(input)
		if err != nil {
			return "", fmt.Errorf("absolute path does not exist: %s", input)
		}
		if !info.IsDir() {
			return "", fmt.Errorf("path is not a directory: %s", input)
		}
		return filepath.Abs(input)
	}

	// 1. Try relative to executable directory
	absFromExe := filepath.Join(exeDir, input)
	if info, err := os.Stat(absFromExe); err == nil && info.IsDir() {
		return filepath.Abs(absFromExe)
	}

	// 2. Try relative to current working directory
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	absFromCwd := filepath.Join(cwd, input)
	if info, err := os.Stat(absFromCwd); err == nil && info.IsDir() {
		return filepath.Abs(absFromCwd)
	}

	return "", fmt.Errorf("directory not found near executable or in current directory: %s", input)
}

func main() {
	// Print banner
	fmt.Println("\n" + colorCyan + banner + colorReset)
	fmt.Println("\nhttps://github.com/MehdiShekari/prompt-generator-for-ai")
	fmt.Println()

	// Determine executable directory (for config file location)
	exePath, err := os.Executable()
	if err != nil {
		exePath = os.Args[0]
	}
	exeDir := filepath.Dir(exePath)

	// Command line flags
	var (
		outputFile string
		configFile string
	)
	flag.StringVar(&outputFile, "output", "", "Output text file (default: project_prompt_<timestamp>.txt)")
	flag.StringVar(&configFile, "config", filepath.Join(exeDir, "config.json"), "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := loadConfig(configFile)
	if err != nil {
		fmt.Printf("%sError loading config: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}

	// Ask user for target directory
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the target directory path (press Enter for current dir): ")
	dirInput, _ := reader.ReadString('\n')
	dirInput = strings.TrimSpace(dirInput)

	targetDir, err := resolveTargetDir(dirInput, exeDir)
	if err != nil {
		fmt.Printf("%sError: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}

	fmt.Printf("Scanning directory: %s\n", targetDir)
	structure, contents, err := buildFileTree(targetDir, cfg)
	if err != nil {
		fmt.Printf("%sError scanning: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}

	// Generate output file name
	if outputFile == "" {
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		outputFile = fmt.Sprintf("project_prompt_%s.txt", timestamp)
	}

	// Write output
	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("%sError creating output file: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}
	defer out.Close()

	out.WriteString("=== Project Folder Structure ===\n\n")
	out.WriteString(strings.Join(structure, "\n"))
	out.WriteString("\n\n=== File Contents ===\n")
	out.WriteString(strings.Join(contents, "\n"))

	fmt.Printf("%sExport complete! File saved as: %s%s\n", colorGreen, outputFile, colorReset)
}
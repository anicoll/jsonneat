package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/anicoll/jsonneat/sorter"
)

func main() {
	inPlace := flag.Bool("w", false, "write result to source file instead of stdout")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: jsonneat [-w] <path> [paths...]")
		fmt.Fprintln(os.Stderr, "  -w       write result to source file instead of stdout")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Path can be:")
		fmt.Fprintln(os.Stderr, "  - a specific file (e.g., file.jsonnet)")
		fmt.Fprintln(os.Stderr, "  - a directory (e.g., ./configs)")
		fmt.Fprintln(os.Stderr, "  - ./... for recursive search from current directory")
		os.Exit(1)
	}

	var filesToProcess []string
	for _, arg := range args {
		files, err := expandPath(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error expanding path %s: %v\n", arg, err)
			os.Exit(1)
		}
		filesToProcess = append(filesToProcess, files...)
	}

	if len(filesToProcess) == 0 {
		fmt.Fprintln(os.Stderr, "No jsonnet files found")
		os.Exit(1)
	}

	for _, filePath := range filesToProcess {
		if err := processFile(filePath, *inPlace); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", filePath, err)
			os.Exit(1)
		}
	}
}

// expandPath expands a path argument into a list of jsonnet files
func expandPath(path string) ([]string, error) {
	// Handle ./... pattern for recursive search
	if path == "./..." {
		return findJsonnetFiles(".", true)
	}

	// Check if path ends with /... for recursive search
	if filepath.Base(path) == "..." {
		dir := filepath.Dir(path)
		return findJsonnetFiles(dir, true)
	}

	// Check if it's a file or directory
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		// Non-recursive directory search
		return findJsonnetFiles(path, false)
	}

	// Single file
	if filepath.Ext(path) != ".jsonnet" && filepath.Ext(path) != ".libsonnet" {
		return nil, fmt.Errorf("file must have .jsonnet or .libsonnet extension")
	}
	return []string{path}, nil
}

// findJsonnetFiles finds all jsonnet files in a directory
func findJsonnetFiles(dir string, recursive bool) ([]string, error) {
	var files []string

	if recursive {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && (filepath.Ext(path) == ".jsonnet" || filepath.Ext(path) == ".libsonnet") {
				files = append(files, path)
			}
			return nil
		})
		return files, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			if filepath.Ext(name) == ".jsonnet" || filepath.Ext(name) == ".libsonnet" {
				files = append(files, filepath.Join(dir, name))
			}
		}
	}

	return files, nil
}

func processFile(filePath string, inPlace bool) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	sorted, err := sorter.SortJsonnet(string(content))
	if err != nil {
		return fmt.Errorf("sorting content: %w", err)
	}

	// Cleanup whitespace and formatting
	cleaned := sorter.CleanupWhitespace(sorted)

	if inPlace {
		if err := os.WriteFile(filePath, []byte(cleaned), 0644); err != nil {
			return fmt.Errorf("writing file: %w", err)
		}
		fmt.Printf("Formatted %s\n", filePath)
	} else {
		fmt.Print(cleaned)
	}

	return nil
}

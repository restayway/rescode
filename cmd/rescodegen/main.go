// Package main provides the rescodegen command-line tool for generating
// type-safe Go error code constants and factory functions from YAML/JSON definitions.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/restayway/rescode/internal/generator"
)

const version = "1.0.0"

func main() {
	var (
		input   = flag.String("input", "", "Path to YAML/JSON file containing error definitions (required)")
		output  = flag.String("output", "rescode_gen.go", "Path to generated Go file")
		pkg     = flag.String("package", "", "Go package name to use in generated code (defaults to package of output file directory)")
		showVer = flag.Bool("version", false, "Show version information")
		help    = flag.Bool("help", false, "Show help information")
	)
	
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	if *showVer {
		fmt.Printf("rescodegen version %s\n", version)
		return
	}

	if *input == "" {
		fmt.Fprintf(os.Stderr, "Error: --input is required\n\n")
		showHelp()
		os.Exit(1)
	}

	// Open input file
	inputFile, err := os.Open(*input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to open input file %s: %v\n", *input, err)
		os.Exit(1)
	}
	defer inputFile.Close()

	// Parse error definitions
	errors, err := generator.ParseInput(inputFile, *input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to parse input file: %v\n", err)
		os.Exit(1)
	}

	// Determine package name
	packageName := *pkg
	if packageName == "" {
		// Default to the directory name of the output file
		dir := filepath.Dir(*output)
		if dir == "." {
			dir, _ = os.Getwd()
		}
		packageName = filepath.Base(dir)
	}

	// Generate code
	config := generator.Config{
		Package: packageName,
		Errors:  errors,
	}
	
	code, err := generator.Generate(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to generate code: %v\n", err)
		os.Exit(1)
	}

	// Write output file
	if err := os.WriteFile(*output, code, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to write output file %s: %v\n", *output, err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated %s with %d error definitions\n", *output, len(errors))
}

func showHelp() {
	fmt.Printf(`rescodegen - Type-Safe Go Error Code Generator

Usage:
  rescodegen --input <file> [--output <file>] [--package <name>]

Options:
  --input     Path to YAML/JSON file containing error definitions (required)
  --output    Path to generated Go file (default: rescode_gen.go)
  --package   Go package name to use in generated code (default: directory name)
  --version   Show version information
  --help      Show this help message

Examples:
  rescodegen --input errors.yaml --output rescode_gen.go --package myservice
  go run github.com/restayway/rescode/cmd/rescodegen --input errors.json

For go:generate usage:
  //go:generate go run github.com/restayway/rescode/cmd/rescodegen --input errors.yaml --output rescode_gen.go --package myservice

Input file format (YAML):
  - code: 20001
    key: PolicyNotFound
    message: Policy not found
    http: 404
    grpc: 5
    desc: Policy could not be located in the database

Input file format (JSON):
  [
    {
      "code": 20001,
      "key": "PolicyNotFound", 
      "message": "Policy not found",
      "http": 404,
      "grpc": 5,
      "desc": "Policy could not be located in the database"
    }
  ]
`)
}
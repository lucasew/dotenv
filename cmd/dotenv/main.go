package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
)

var env = map[string]string{}

// mergeEnv takes a parsed dotenv map and merges it into the global env state.
// This ensures that subsequent variables can override previous ones.
func mergeEnv(m map[string]string) {
	for k := range m {
		env[k] = m[k]
	}
}

// printHelp displays the usage instructions for the CLI.
// It outlines the two main input formats: @file for .env files and --key=value for inline vars.
func printHelp() {
	println("dotenv [...params] -- command")
	println("params: ")
	println(" @file.txt load file as dotenv")
	println(" --key=value load variable")
}

// parseEnvTerm parses a single command-line argument to extract environment variables.
// It handles two formats:
// 1. "@filename" - Reads the specified file using godotenv and merges it into the global state.
// 2. "--key=value" - Directly injects a key-value pair into the global state.
func parseEnvTerm(term string) error {
	if term[0] == '@' {
		filename := term[1:]
		f, err := os.Open(filename)
		defer f.Close()
		if err != nil {
			return err
		}
		variables, err := godotenv.Parse(f)
		if err != nil {
			return err
		}
		mergeEnv(variables)
		return nil
	}
	if strings.HasPrefix(term, "--") {
		termBody := term[2:]
		elems := strings.Split(termBody, "=")
		if len(elems) != 2 {
			return fmt.Errorf("syntax error near %s", term)
		}
		key := elems[0]
		value := elems[1]
		env[key] = value
		return nil
	}
	fmt.Printf("warn: ignoring argument '%s'\n", term)
	return nil
}

// handleError is the centralized error reporting function for the project.
// It intercepts any non-nil error, prints it to stdout, displays the help menu,
// and strictly terminates the process with exit code 1 to avoid silent failures.
func handleError(err error) {
	if err == nil {
		return
	}
	fmt.Println("error: ", err)
	printHelp()
	os.Exit(1)
}

// main is the entrypoint. It splits os.Args into two logical parts using "--" as a delimiter:
// 1. The environment definition terms (parsed into the global env map).
// 2. The target command and its arguments.
// Once the environment is fully loaded (including an implicit @.env override at the end),
// it spawns the target command as a subprocess, inherits os.Environ(), and appends the custom env map.
func main() {
	argslen := len(os.Args)
	stripped := make([]string, argslen-1)
	for i := 1; i < argslen; i++ {
		stripped[i-1] = strings.Trim(os.Args[i], " ")
	}
	foundDivider := false
	command := []string{}
	for i := 0; i < len(stripped); i++ {
		if stripped[i] == "--" {
			if !foundDivider {
				foundDivider = true
				continue
			}
		}
		if !foundDivider {
			err := parseEnvTerm(stripped[i])
			handleError(err)
		} else {
			command = append(command, stripped[i])
		}
	}
	if !foundDivider {
		handleError(fmt.Errorf("missing divider (--)"))
	}
	parseEnvTerm("@.env")
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Env = os.Environ() // herdar env do pai
	for k, v := range env {
		cmd.Env = append(cmd.Env, strings.Join([]string{k, v}, "="))
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	handleError(err)
	proc, err := cmd.Process.Wait()
	handleError(err)
	os.Exit(proc.ExitCode())
}

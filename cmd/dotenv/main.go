// Package main implements the dotenv CLI tool which loads environment variables
// from files and command-line arguments before executing a command.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
)

// env stores the loaded environment variables.
var env  = map[string]string{}

// mergeEnv updates the global env map with variables from the provided map.
func mergeEnv(m map[string]string) {
    for k := range m {
        env[k] = m[k]
    }
}

// printHelp displays the usage information for the dotenv command.
func printHelp() {
    println("dotenv [...params] -- command")
    println("params: ")
    println(" @file.txt load file as dotenv")
    println(" --key=value load variable")
}

// parseEnvTerm parses a single command-line argument to load environment variables.
//
// It supports two formats:
//   - @filename: Loads environment variables from the specified file.
//   - --key=value: Sets a specific environment variable.
//
// Note: The current implementation of --key=value does not support values containing '='.
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

// handleError prints the error message and help text, then exits the program with status code 1 if an error occurs.
func handleError(err error) {
    if err == nil {
        return
    }
    fmt.Println("error: ", err)
    printHelp()
    os.Exit(1)
}

// main is the entry point of the application.
//
// It parses command-line arguments to load environment variables,
// loads the default .env file (which overrides CLI arguments),
// and finally executes the specified command with the populated environment.
func main() {
    argslen := len(os.Args)
    stripped := make([]string, argslen - 1)
    for i := 1; i < argslen; i++ {
        stripped[i - 1] = strings.Trim(os.Args[i], " ")
    }
    foundDivider := false
    command := []string{}
    for i := 0; i < len(stripped); i++ {
        if (stripped[i] == "--") {
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
    // Load .env file. This overrides any conflicting variables set via CLI arguments.
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

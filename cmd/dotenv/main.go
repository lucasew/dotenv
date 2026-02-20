package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
)

var env = map[string]string{}

func mergeEnv(m map[string]string) {
	for k := range m {
		env[k] = m[k]
	}
}

func printHelp() {
	println("dotenv [...params] -- command")
	println("params: ")
	println(" @file.txt load file as dotenv")
	println(" --key=value load variable")
}

func parseEnvTerm(term string) error {
	if term[0] == '@' {
		filename := term[1:]
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer func() { _ = f.Close() }()

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

func handleError(err error) {
	if err == nil {
		return
	}
	fmt.Println("error: ", err)
	printHelp()
	os.Exit(1)
}

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
	_ = parseEnvTerm("@.env")
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

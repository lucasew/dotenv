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
	if len(term) == 0 {
		return nil
	}
	if term[0] == '@' {
		filename := term[1:]
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
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
	stripped := make([]string, 0, len(os.Args)-1)
	for _, arg := range os.Args[1:] {
		stripped = append(stripped, strings.Trim(arg, " "))
	}
	foundDivider := false
	command := []string{}
	for _, arg := range stripped {
		if arg == "--" {
			if !foundDivider {
				foundDivider = true
				continue
			}
		}
		if !foundDivider {
			err := parseEnvTerm(arg)
			handleError(err)
		} else {
			command = append(command, arg)
		}
	}
	if !foundDivider {
		handleError(fmt.Errorf("missing divider (--)"))
	}
	if len(command) == 0 {
		handleError(fmt.Errorf("missing command"))
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

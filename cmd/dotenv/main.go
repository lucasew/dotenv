package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func printHelp() {
	println("dotenv [...params] -- command")
	println("params: ")
	println(" @file.txt load file as dotenv")
	println(" --key=value load variable")
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
	env := map[string]string{}
	for i := 0; i < len(stripped); i++ {
		if stripped[i] == "--" {
			if !foundDivider {
				foundDivider = true
				continue
			}
		}
		if !foundDivider {
			err := ParseEnvTerm(env, stripped[i])
			handleError(err)
		} else {
			command = append(command, stripped[i])
		}
	}
	if !foundDivider {
		handleError(fmt.Errorf("missing divider (--)"))
	}
	_ = ParseEnvTerm(env, "@.env")
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

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func mergeEnv(env map[string]string, m map[string]string) {
	for k := range m {
		env[k] = m[k]
	}
}

func parseEnvTerm(term string, env map[string]string) error {
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
		mergeEnv(env, variables)
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

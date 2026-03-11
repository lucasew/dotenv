package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func ParseEnvTerm(env map[string]string, term string) error {
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
		for k, v := range variables {
			env[k] = v
		}
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

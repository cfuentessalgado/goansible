package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Go Ansible")
	histfile := os.Getenv("HOME")
	fmt.Println(histfile)
	f, err := os.Open(histfile + "/.zsh_history")
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)
	installs := Installs{}
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ":")
		if len(line) > 2 && strings.Contains(line[2], "sudo dnf install") {
			installs.Commands = append(installs.Commands, line[2])
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	successes := Installs{}
	for _, c := range installs.Commands {
		parts := strings.Split(c, ";")
		status, _ := strconv.Atoi(parts[0])
		command := parts[1]
		if status == 0 {
			successes.Commands = append(successes.Commands, command)
		}
	}

	suggestions := prepareSuggestions(successes.Commands)
	cleanSuggestions(&suggestions)
	fmt.Println(suggestions)
}

func prepareSuggestions(commands []string) Installs {
	suggestions := Installs{}
	for _, c := range commands {
		install := strings.Split(strings.Replace(c, "sudo dnf install ", "", 1), " ")

		suggestions.Commands = append(suggestions.Commands, install...)
	}
	return suggestions
}

func cleanSuggestions(suggestions *Installs) {
    clean := Installs{}
    for _, c := range suggestions.Commands {
        c = strings.Replace(c, "\\", "",1)
        if len(c) > 3 {
            clean.Commands = append(clean.Commands, c)
        }
    }
    suggestions.Commands = clean.Commands
}

type Installs struct {
	Commands []string
}

package env

import "github.com/urfave/cli"

var commands []cli.Command

// Command register commands
func Command(args ...cli.Command) {
	commands = append(commands, args...)
}

package main

import (
	"log"

	"github.com/mitchellh/cli"
)

type ListCommand struct {
}

func (c *ListCommand) Help() string {
	return "List all shares with published results."
}

func (c *ListCommand) Synopsis() string {
	return c.Help()
}

func (c *ListCommand) Run(args []string) int {
	ipos, err := getIPOList()

	if err != nil {
		log.Println(err)
		return 1
	}

	printIPOTable(ipos)

	return 0
}

func ListCommandFactory() (cli.Command, error) {
	return &ListCommand{}, nil
}

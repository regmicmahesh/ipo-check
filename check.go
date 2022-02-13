package main

import (
	"fmt"
	"log"

	"github.com/mitchellh/cli"
)

type CheckCommand struct {
}

func (c *CheckCommand) Help() string {
	return "Check multiple share by passing boid separated by space."
}

func (c *CheckCommand) Synopsis() string {
	return c.Help()
}

func (c *CheckCommand) Run(args []string) int {
	ipos, err := getIPOList()

	if err != nil {
		log.Println(err)
		return 1
	}

	if len(args) == 0 {
		fmt.Println("Usage: check <boid1> <boid2> ...")
		return 1
	}

	ipoMap := make(map[int]IPOInfo)
	for _, ipo := range ipos {
		ipoMap[ipo.ID] = ipo
	}

	type ShareStatus struct {
		BOID    string
		success bool
		ID      int
	}

	success := make(chan ShareStatus)

	for _, ipo := range ipos {

		for _, boid := range args {
			go func(ipo IPOInfo, boid string) {

				status, _ := checkIPO(boid, ipo.ID)

				success <- ShareStatus{
					success: status,
					ID:      ipo.ID,
					BOID:    boid,
				}
			}(ipo, boid)

		}

	}

	for i := 0; i < len(ipos)*len(args); i++ {
		status := <-success
		if status.success {
			fmt.Printf("%s: %s\n", status.BOID, ipoMap[status.ID].Name)
		}
	}

	return 0

}

func CheckCommandFactory() (cli.Command, error) {
	return &CheckCommand{}, nil
}

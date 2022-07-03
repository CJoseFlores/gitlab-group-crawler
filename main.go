package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	var args ProgArgs

	app := &cli.App{
		Name:  "gitlab-project-name-crawler",
		Usage: "Prints out the names of projects underneath a GitLab group (recursively)",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "insecure",
				Value:       false,
				Usage:       "connect to gitlab without TLS verification",
				Aliases:     []string{"k"},
				Destination: &args.Insecure,
			},
		},
		Action: func(ctx *cli.Context) error {
			fmt.Println("Unimplemented")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

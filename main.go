package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/xanzy/go-gitlab"
)

func main() {
	var args ProgArgs

	app := &cli.App{
		Name:  "gitlab-project-name-crawler",
		Usage: "Prints out the names of projects underneath a GitLab group (recursively)",
		UsageText: "[global options] command [command options] [arguments...]" +
			"\n\nThe list of arguments passed in are gitlab groups or subgroups you wish to scan recursively",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "insecure",
				Value:       false,
				Usage:       "connect to gitlab without TLS verification",
				Aliases:     []string{"k"},
				Destination: &args.Insecure,
			},
			&cli.StringFlag{
				Name:        "file-name",
				Usage:       "the name of an existing file to diff against",
				Aliases:     []string{"f"},
				Destination: &args.InputFileName,
			},
			&cli.StringFlag{
				Name:        "output-file-name",
				Value:       "project-list.txt",
				Usage:       "the name of the output file to generate",
				Aliases:     []string{"o"},
				Destination: &args.OutputFileName,
			},
			&cli.StringFlag{
				Name:        "gitlab-url",
				Value:       "https://gitlab.com/",
				Usage:       "the URL of the gitlab instance to crawl",
				Aliases:     []string{"g"},
				Destination: &args.GitlabUrl,
			},
			&cli.StringFlag{
				Name:        "username",
				Required:    true,
				Usage:       "the name of the account that has access to the groups to scan (REQUIRED)",
				Aliases:     []string{"u"},
				Destination: &args.Username,
			},
			&cli.StringFlag{
				Name:        "password",
				Required:    true,
				Usage:       "the password of the account that has access to the groups to scan (REQUIRED)",
				Aliases:     []string{"p"},
				Destination: &args.Password,
			},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() != 1 {
				return errors.New("the crawler currently only supports scanning 1 group")
			}

			fmt.Println("Argument: " + ctx.Args().First())
			fmt.Println("Unimplemented")
			return nil
		},
		EnableBashCompletion: true,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	git, err := gitlab.NewBasicAuthClient(args.Username, args.Password, gitlab.WithBaseURL("https://gitlab.com"))
	fmt.Print(git.BaseURL(), err)
	fmt.Println()
}

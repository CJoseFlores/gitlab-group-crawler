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
	args := parseArgs()

	// Exit early if no arguments were passed in
	if len(args.Groups) <= 0 {
		os.Exit(0)
	}

	git, err := gitlab.NewClient(args.GitlabToken, gitlab.WithBaseURL(args.GitlabUrl))
	if err != nil {
		log.Fatal(err)
	}

	projects, response, err := git.Groups.ListGroupProjects(
		// FIXME: Hard-coded only ever looking at the first group
		args.Groups[0],
		&gitlab.ListGroupProjectsOptions{
			Archived:         gitlab.Bool(false),
			IncludeSubGroups: gitlab.Bool(true),
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	if response.StatusCode != 200 {
		log.Fatalf("Could not fetch groups (Code: %v)\n", response.StatusCode)
	}

	for _, project := range projects {
		fmt.Println(project.PathWithNamespace)
	}
}

func parseArgs() ProgArgs {
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
				Name:        "gitlab-token",
				Required:    true,
				Usage:       "a token of the account that has read access to the groups to scan (REQUIRED)",
				Aliases:     []string{"t"},
				Destination: &args.GitlabToken,
			},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() < 1 {
				return errors.New("no groups provided")
			}
			args.Groups = ctx.Args().Slice()
			return nil
		},
		EnableBashCompletion: true,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	return args
}

package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/xanzy/go-gitlab"
)

// Holds the version of the program.
// Should be set upon building via '-ldflags "-X main.version=VERSION"'
var version string

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

	scanGroups(git, args)
}

// Parses program arguments.
func parseArgs() ProgArgs {
	var args ProgArgs

	app := &cli.App{
		Name:  "gitlab-project-name-crawler",
		Usage: "Prints out the names of projects underneath a GitLab group (recursively)",
		UsageText: "[global options] command [command options] [arguments...]" +
			"\n\nThe list of arguments passed in are gitlab groups or subgroups you wish to scan recursively",
		Version: version,
		Flags: []cli.Flag{
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

// Scans the requested Gitlab groups and outputs a file with their paths.
func scanGroups(git *gitlab.Client, args ProgArgs) {
	outputFile, err := os.Create(args.OutputFileName)
	if err != nil {
		fmt.Println("Could not create output file...", err)
	}
	defer outputFile.Close()

	// List the projects for every specified group
	for _, group := range args.Groups {

		page := 1
		for {
			projects, response, err := git.Groups.ListGroupProjects(
				group,
				&gitlab.ListGroupProjectsOptions{
					Archived:         gitlab.Bool(false),
					IncludeSubGroups: gitlab.Bool(true),
					ListOptions:      gitlab.ListOptions{PerPage: 100, Page: page},
				},
			)

			if err != nil {
				log.Fatal(err)
			}
			if response.StatusCode != 200 {
				log.Fatalf("Could not fetch groups (Code: %v)\n", response.StatusCode)
			}

			// Print out the list of discovered projects and write them to file
			for _, project := range projects {
				fmt.Println(project.PathWithNamespace)
				outputFile.WriteString(project.PathWithNamespace + "\n")
			}

			// Continue onto next page, or exit if we are on the last page
			if page < response.TotalPages {
				page = response.NextPage
			} else {
				break
			}
		}
	}
}

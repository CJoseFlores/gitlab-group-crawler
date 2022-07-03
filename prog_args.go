package main

// Contains the arguments passed into the program
type ProgArgs struct {
	// Whether or not to verify Gitlab's certificate
	Insecure bool
	// The name of the file to diff against
	InputFileName string
	// The name of the file to output
	OutputFileName string
	// The URL of the gitlab instance (ex: https://gitlab.com)
	GitlabUrl string
	// The token of the account that has read access to the groups to scan
	GitlabToken string
	// The list of groups to scan
	Groups []string
}

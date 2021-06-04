package ymir

// import (
// 	"fmt"

// 	"github.com/google/go-github/github"
// 	"golang.org/x/oauth2"
// )

// // Sudo plan
// func testGitAuth(cmd YmirCommand) error {
// 	ts := oauth2.StaticTokenSource(
// 		&oauth2.Token{AccessToken: cmd.App.GetConfig().Git.Github.AccessToken},
// 	)
// 	ctx := cmd.Context.Context

// 	tc := oauth2.NewClient(ctx, ts)

// 	client := github.NewClient(tc)

// 	opt := github.RepositoryListOptions{
// 		ListOptions: github.ListOptions{
// 			PerPage: 1,
// 		},
// 	}

// 	_, _, err := client.Repositories.List(cmd.Context.Context, "", &opt)

// 	if err != nil {
// 		fmt.Printf("%#v\n", err)
// 		return err
// 	}

// 	fmt.Print("Successfully connected to Github API!\n")
// 	return nil
// }

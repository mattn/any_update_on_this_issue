package main

import (
	"code.google.com/p/goauth2/oauth"
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"log"
	"net/http"
	"os"
	"strconv"
)

var token = flag.String("token", "", "Token to access github API")

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage: %s [user] [repo] [issue-number]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.Parse()
	if flag.NArg() != 3 {
		flag.Usage()
	}
	if token == nil || *token == "" {
		flag.Usage()
	}
	issue, err := strconv.Atoi(flag.Arg(2))
	if err != nil {
		flag.Usage()
	}
	var httpClient *http.Client
	var client *github.Client
	if *token != "" {
		httpClient = (&oauth.Transport{
			Token: &oauth.Token{AccessToken: *token},
		}).Client()
	}
	client = github.NewClient(httpClient)
	res, _, err := client.Issues.Get(flag.Arg(0), flag.Arg(1), issue)
	if err != nil {
		log.Fatal(err)
	}
	if res.State != nil && *res.State != "open" {
		return
	}

	message := "Any update on this issue?"
	var comment github.IssueComment
	comment.Body = &message
	_, _, err = client.Issues.CreateComment(flag.Arg(0), flag.Arg(1), issue, &comment)
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"context"
	"flag"
	"github.com/whosonfirst/go-whosonfirst-readwrite-github/reader"
	"io"
	"log"
	"os"
)

func main() {

	var owner = flag.String("owner", "whosonfirst-data", "...")
	var repo = flag.String("repo", "whosonfirst-data", "...")
	var branch = flag.String("branch", "master", "...")
	var token = flag.String("token", "", "...")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r, err := reader.NewGitHubAPIReader(ctx, *owner, *repo, *branch, *token)

	if err != nil {
		log.Fatal(err)
	}

	for _, path := range flag.Args() {

		fh, err := r.Read(path)

		if err != nil {
			log.Fatal(err)
		}

		io.Copy(os.Stdout, fh)
	}

}

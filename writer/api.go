package writer

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	wof_writer "github.com/whosonfirst/go-whosonfirst-readwrite/writer"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	_ "log"
	"time"
)

type GitHubAPIWriter struct {
	wof_writer.Writer
	owner    string
	repo     string
	branch   string
	client   *github.Client
	context  context.Context
	throttle <-chan time.Time
}

func NewGitHubAPIWriter(ctx context.Context, owner string, repo string, branch string, token string) (wof_writer.Writer, error) {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// https://github.com/golang/go/wiki/RateLimiting

	rate := time.Second / 3
	throttle := time.Tick(rate)

	r := GitHubAPIWriter{
		repo:     repo,
		owner:    owner,
		branch:   branch,
		throttle: throttle,
		client:   client,
		context:  ctx,
	}

	return &r, nil
}

func (r *GitHubAPIWriter) Write(path string, fh io.ReadCloser) error {

	<-r.throttle

	body, err := ioutil.ReadAll(fh)

	if err != nil {
		return err
	}

	url := r.URI(path)

	commit_msg := ""
	name := ""
	email := ""

	update_opts := &github.RepositoryContentFileOptions{
		Message: github.String(commit_msg),
		Content: body,
		Branch:  github.String(r.branch),
		Committer: &github.CommitAuthor{
			Name:  github.String(name),
			Email: github.String(email),
		},
	}

	get_opts := &github.RepositoryContentGetOptions{}

	get_rsp, _, _, err := r.client.Repositories.GetContents(r.context, r.owner, r.repo, url, get_opts)

	if err != nil {
		update_opts.SHA = get_rsp.SHA
	}

	_, _, err = r.client.Repositories.UpdateFile(r.context, r.owner, r.owner, url, update_opts)

	if err != nil {
		return nil
	}

	return nil
}

func (r *GitHubAPIWriter) URI(key string) string {

	return fmt.Sprintf("data/%s", key)
}

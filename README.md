# go-whosonfirst-readwrite-github

## Important

This package is officially deprecated and has been superseded by the [go-reader-*](https://github.com/whosonfirst?utf8=%E2%9C%93&q=go-reader&type=&language=), [go-writer-*](https://github.com/whosonfirst?utf8=%E2%9C%93&q=go-writer&type=&language=) packages, as well as [go-whosonfirst-reader](https://github.com/whosonfirst/go-whosonfirst-reader) and [go-whosonfirst-writer](https://github.com/whosonfirst/go-whosonfirst-writer). In time this repository will be archived.

## Install

You will need to have both `Go` (specifically version [1.12](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make tools
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Important

This is work in progress and should be considered to work "until it doesn't". Things may change still.

## Example

### Deprecating one or more Who's On First records (using the GitHub API)

_Note the use of the `github.com/tidwall/sjson` package which is not part of this package._

```
package main

import (
	"bytes"
	"context"
	"flag"
	"github.com/tidwall/sjson"
	"github.com/whosonfirst/go-whosonfirst-readwrite-github/reader"
	"github.com/whosonfirst/go-whosonfirst-readwrite-github/writer"
	"io"
	"io/ioutil"
	"log"
	"time"
)

// really this function should be part of the go-whosonfirst-export package
// but for the purposes of example code it will do...

func deprecate(fh io.ReadCloser) (io.ReadCloser, error) {

	body, err := ioutil.ReadAll(fh)

	if err != nil {
		return nil, err
	}

	now := time.Now()

	body, err = sjson.SetBytes(body, "properties.edtf:deprecated", now.Format("2006-01-02"))

	if err != nil {
		return nil, err
	}

	body, err = sjson.SetBytes(body, "properties.mz:is_current", "0")

	if err != nil {
		return nil, err
	}

	body, err = sjson.SetBytes(body, "wof:lastmodified", now.Unix())

	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(body)
	out := ioutil.NopCloser(r)

	return out, nil
}

func main() {

	var owner = flag.String("owner", "whosonfirst-data", "A valid GitHub user or organizartion")
	var repo = flag.String("repo", "whosonfirst-data", "A valid GitHub repository")
	var branch = flag.String("branch", "master", "A valid Git branch")
	var token = flag.String("token", "", "A valid GitHub API token")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r, err := reader.NewGitHubAPIReader(ctx, *owner, *repo, *branch, *token)

	if err != nil {
		log.Fatal(err)
	}

	templates := &writer.GitHubAPIWriterCommitTemplates{
		New: "Initial commit",
		Update: "Deprecate %s",
	}
	
	wr, err := writer.NewGitHubAPIWriter(ctx, *owner, *repo, *branch, *token, templates)

	if err != nil {
		log.Fatal(err)
	}

	for _, path := range flag.Args() {

		fh, err := r.Read(path)

		if err != nil {
			log.Fatal(err)
		}

		fh2, err := deprecate(fh)

		if err != nil {
			log.Fatal(err)
		}

		err = wr.Write(path, fh2)

		if err != nil {
			log.Fatal(err)
		}
	}

}
```

## See also

* https://github.com/whosonfirst/go-whosonfirst-readwrite
* https://github.com/google/go-github

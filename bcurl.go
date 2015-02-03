// Copyright 2014 Fabio Rehm. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/fgrehm/boom-curl/boom/boomer"
)

const (
	headerRegexp = "^([\\w-]+):\\s*(.+)"
	authRegexp   = "^([\\w-\\.]+):(.+)"
)

func main() {
	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.Name}} {{if .Flags}}[options] {{end}}<URL>
VERSION:
   {{.Version}}{{if or .Author .Email}}
AUTHOR:{{if .Author}}
  {{.Author}}{{if .Email}} - <{{.Email}}>{{end}}{{else}}
  {{.Email}}{{end}}{{end}}
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}
`

	app := cli.NewApp()
	app.Name = "bcurl"
	app.Usage = "A cURL like interface for Boom, an HTTP(S) load generator, ApacheBench (ab) replacement"
	app.Author = "Fábio Rehm"
	app.Email = "fgrehm@gmail.com"
	app.Version = VERSION
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "H, header",
			Value: &cli.StringSlice{},
			Usage: "custom header to pass to server",
		},
		cli.StringFlag{
			Name:  "d, data",
			Value: "",
			Usage: "HTTP POST data",
		},
		cli.IntFlag{
			Name:  "cpus",
			Value: runtime.NumCPU(),
			Usage: fmt.Sprintf("Number of used cpu cores. (default for current machine is %d cores)", runtime.NumCPU()),
		},
		cli.IntFlag{
			Name:  "requests",
			Value: 200,
			Usage: "Number of requests to run",
		},
		cli.IntFlag{
			Name:  "concurrency",
			Value: 50,
			Usage: "Number of requests to run concurrently.",
		},
	}
	app.Action = func(c *cli.Context) {
		if len(c.Args()) != 1 {
			usageAndExit(c, "")
		}

		boom(c)
	}

	// cURL accepts the URL as the first argument but the cli package we are
	// using will blow up since it will recognize it as a subcommand that does
	// not exist, that's why we need to move things around
	if len(os.Args) > 2 && !strings.HasPrefix(os.Args[1], "-") {
		argsWithoutUrl := append([]string{os.Args[0]}, os.Args[2:]...)
		os.Args = append(argsWithoutUrl, os.Args[1])
	}
	app.Run(os.Args)
}

func boom(c *cli.Context) {
	body := c.String("data")
	var m string
	if body != "" {
		m = "POST"
	} else {
		m = "GET"
	}
	cpus := c.Int("cpus")
	num := c.Int("requests")
	conc := c.Int("concurrency")

	q := 0
	t := 0
	contentType := "text/html"
	if m == "POST" {
		contentType = "application/x-www-form-urlencoded"
	}
	output := ""
	insecure := false

	// MAGIC HAPPENS BELOW
	runtime.GOMAXPROCS(cpus)

	if num <= 0 || conc <= 0 {
		usageAndExit(c, "n and c cannot be smaller than 1.")
	}

	var (
		url, method string
		header http.Header = make(http.Header)
	)

	method = strings.ToUpper(m)
	url = c.Args()[0]

	// Prepend http protocol if user did not provide it
	if !regexp.MustCompile(`^https?://`).MatchString(url) {
		url = "http://" + url
	}

	// set content-type
	header.Set("User-Agent", "boom-curl/"+VERSION)
	header.Set("Content-Type", contentType)
	for _, h := range c.StringSlice("H") {
		headerAndValue := strings.SplitAfterN(h, ":", 2)
		h := strings.TrimSuffix(headerAndValue[0], ":")
		value := strings.Trim(strings.Trim(headerAndValue[1], " "), "\"")

		header.Set(h, value)
	}

	if output != "csv" && output != "" {
		usageAndExit(c, "Invalid output type.")
	}

	(&boomer.Boomer{
		Req: &boomer.ReqOpts{
			Method: method,
			URL:    url,
			Body:   body,
			Header: header,
		},
		N:             num,
		C:             conc,
		Qps:           q,
		Timeout:       t,
		AllowInsecure: insecure,
		Output:        output,
	}).Run()
}

func usageAndExit(c *cli.Context, message string) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	cli.ShowAppHelp(c)
	os.Exit(1)
}

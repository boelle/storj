// Copyright (C) 2024 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
)

func main() {
	args, err := parseFlags()
	if err != nil {
		exitWith(100, "invalid command line parameter. %s", err)
	}
}

type params struct {
	port       uint
	backendURL *url.URL
}

func parseFlags() (params, error) {
	fl := flag.NewFlagSet("", flag.ExitOnError)
	fl.Usage = func() {
		fmt.Printf(
			"Run a proxy for satellite admin Web App development to be able to test different user's roles.\n\n",
		)
		fl.PrintDefaults()
	}

	var (
		args = params{}
		err  error
	)

	port := flag.Uint("port", 0, "the port where this proxy will listen")
	backendURL := flag.String("backend-url", "", "absolute URL to the backend to proxy (e.g. http://localhost:10005)")

	flag.Parse()

	args.port = *port
	args.backendURL, err = url.Parse(*backendURL)
	if err != nil {
		return args, fmt.Errorf("invalid backend URL. %w", err)
	}

	return args, nil
}

// exitWith prints to stderr the message formed by format and a appending a new line and exit the
// program with code.
func exitWith(code int, format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(code)
}

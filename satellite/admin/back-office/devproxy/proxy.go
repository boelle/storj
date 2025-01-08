// Copyright (C) 2024 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/zeebo/errs"

	"storj.io/common/errs2"
)

// Proxy proxies requests to a back-end injecting the X-Forwarded-Groups header for the user
// indicated by the "user" query parameter.
type Proxy struct {
	server http.Server

	// users is a map of users and the groups they belong to.
	users map[string][]string
}

// NewProxy creates a new Proxy instance listing on localhost on the provided port and proxying
// the backend.
func NewProxy(port uint16, backend *url.URL) Proxy {
	return Proxy{
		server: http.Server{
			Addr:    fmt.Sprintf("127.0.0.1:%d", port),
			Handler: httputil.NewSingleHostReverseProxy(backend),
		},
		users: make(map[string][]string),
	}
}

func (p *Proxy) Serve(ctx context.Context) error {
	group := errs2.Group{}

	group.Go(func() error {
		err := p.server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	})

	group.Go(func() error {
		<-ctx.Done()
		return p.server.Shutdown(ctx)
	})

	return errs.Combine(group.Wait()...)
}

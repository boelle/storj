// Copyright (C) 2024 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewProxy(port uint, backend *url.URL) http.Server {
	return http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
		Handler: httputil.NewSingleHostReverseProxy(backend),
	}
}

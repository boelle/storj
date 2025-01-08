// Copyright (C) 2024 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"storj.io/common/process"
)

var rootCmd = &cobra.Command{
	Use:   "devproxy <proxy-port> <backend-url>",
	Short: "Run a proxy for Satellite Admin Web App development to manage authorization",
	Long: `"This command helps to work with the Satellite Admin Web App development.

The Satellite Admin Web App requires a proxy for the users authorization. It is usually a Oauth2
proxy.

Dealing with a Oauth2 proxy on a development environment is tedious. This command runs a simple
proxy that through a simple web interface allows to configure users and roles and through the proxy
you can indicate the user to use with the "user" query parameter.

<proxy-port> is the port where the proxy will listen.
<backend-url> is the URL where the Satellite Admin Web App is listening
(e.g. http://localhost:10005).
`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, _ := process.Ctx(cmd)
		return runProxy(ctx, args[0], args[1])
	},
}

func main() {
	process.ExecWithCustomOptions(rootCmd, process.ExecOptions{
		LoadConfig: func(_ *cobra.Command, _ *viper.Viper) error {
			return nil
		},
		LoggerFactory: func(logger *zap.Logger) *zap.Logger {
			newLogger, level, err := process.NewLogger("placement-test")
			if err != nil {
				panic(err)
			}
			level.SetLevel(zap.DebugLevel)
			return newLogger
		},
	})
}

func runProxy(ctx context.Context, proxyPort string, adminURL string) error {
	port, err := strconv.ParseUint(proxyPort, 10, 16)
	if err != nil {
		return fmt.Errorf("failed to parse proxy port %s. %w", proxyPort, err)
	}

	backend, err := url.Parse(adminURL)
	if err != nil {
		return fmt.Errorf("failed to parse backend URL %s. %w", adminURL, err)
	}

	proxy := NewProxy(uint16(port), backend)
	return proxy.Serve(ctx)
}

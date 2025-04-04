// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2025 Matthew Penner

//go:build linux

package sdlisten_test

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/matthewpi/sd/sdlisten"
)

func Example() {
	// This is just a placeholder context.
	ctx := context.Background()

	// Get all the listeners passed to us by systemd.
	listeners, err := sdlisten.Listeners()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "failed to get listeners from systemd", slog.Any("err", err))
		os.Exit(1)
		return
	}

	// Add a basic handler for `GET /`.
	http.HandleFunc("GET /", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("Hello, world!\n"))
	})

	// Serve an HTTP server on all the listeners.
	for _, l := range listeners {
		// NOTE: while this is the easiest way to Serve a HTTP server for the
		// purposes of this example, you should likely construct your own
		// [http.Server]. Using [http.Serve] doesn't allow you to configure
		// timeouts which can cause a security risk to publicly exposed
		// applications, hence the `nolint` comment.
		_ = http.Serve(l, nil) //nolint:gosec
	}
}

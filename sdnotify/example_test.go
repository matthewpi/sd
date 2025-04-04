// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2025 Matthew Penner

//go:build linux

package sdnotify_test

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/matthewpi/sd/sdnotify"
)

func Example_notify() {
	// Setup a cancelable context.
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	go func() {
		// Set up channel on which to send signal notifications.
		// We must use a buffered channel or risk missing the signal
		// if we're not ready to receive when the signal is sent.
		c := make(chan os.Signal, 1)

		// Passing no signal argument to [signal.Notify] causes
		// all signals to be sent to the channel.
		signal.Notify(c)

		for s := range c {
			switch s {
			case syscall.SIGHUP:
				// Notify systemd that we are reloading the application
				_ = sdnotify.Reloading()

				err := func() error { return nil }() // FIXME: replace with application reload logic
				if err != nil {
					// Notify systemd of the reload failure.
					//
					// NOTE: you can optionally provide your own errno if you
					// have one, if you don't have one, you MUST provide an
					// integer greater than 0.
					_ = sdnotify.Error(err, 1)
					continue
				}

				// No error, notify systemd that the reload was successful.
				_ = sdnotify.Ready()
			case os.Interrupt:
				_ = sdnotify.Stopping()
				stop()
				return
			}
		}
	}()

	// FIXME: setup application, start whatever services or listeners.
	//
	// NOTE: make sure when you are adding your application code, make sure none
	// of it blocks for an extended period of time, a common mistake is starting
	// a HTTP server on the main thread using [http.ListenAndServe] without
	// placing it in a go-routine.

	// Wait until the context is canceled, via a Interrupt signal.
	<-ctx.Done()

	// NOTE: you can add any application cleanup code here, or utilize `defer`
	// calls while starting your application.
}

func Example_watchdog() {
	// This context is just a placeholder, it should be replaced by a context
	// that gets canceled when the application is stopping.
	ctx := context.Background()

	// Get the watchdog interval (if configured).
	i, err := sdnotify.WatchdogInterval()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "failed to get watchdog interval from systemd", slog.Any("err", err))
		os.Exit(1)
		return
	}
	if i > 0 {
		// Send keep-alives to systemd in the background.
		go func(ctx context.Context, i time.Duration) {
			t := time.NewTicker(i)
			defer t.Stop()

			for {
				select {
				case <-ctx.Done():
					break
				case <-t.C:
					if err := sdnotify.Watchdog(); err != nil {
						slog.LogAttrs(ctx, slog.LevelError, "failed to send keep-alive to watchdog", slog.Any("err", err))
					}
				}
			}
		}(ctx, i)
	}
}

func Example_full() {
	// Setup a cancelable context.
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	// Get the watchdog interval (if configured).
	i, err := sdnotify.WatchdogInterval()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "failed to get watchdog interval from systemd", slog.Any("err", err))
		os.Exit(1)
		return
	}
	if i > 0 {
		// Send keep-alives to systemd in the background.
		go func(ctx context.Context, i time.Duration) {
			t := time.NewTicker(i)
			defer t.Stop()

			for {
				select {
				case <-ctx.Done():
					break
				case <-t.C:
					if err := sdnotify.Watchdog(); err != nil {
						slog.LogAttrs(ctx, slog.LevelError, "failed to send keep-alive to watchdog", slog.Any("err", err))
					}
				}
			}
		}(ctx, i)
	}

	go func() {
		// Set up channel on which to send signal notifications.
		// We must use a buffered channel or risk missing the signal
		// if we're not ready to receive when the signal is sent.
		c := make(chan os.Signal, 1)

		// Passing no signal argument to [signal.Notify] causes
		// all signals to be sent to the channel.
		signal.Notify(c)

		for s := range c {
			switch s {
			case syscall.SIGHUP:
				// Notify systemd that we are reloading the application
				_ = sdnotify.Reloading()

				err := func() error { return nil }() // FIXME: replace with application reload logic
				if err != nil {
					// Notify systemd of the reload failure.
					//
					// NOTE: you can optionally provide your own errno if you
					// have one, if you don't have one, you MUST provide an
					// integer greater than 0.
					_ = sdnotify.Error(err, 1)
					continue
				}

				// No error, notify systemd that the reload was successful.
				_ = sdnotify.Ready()
			case os.Interrupt:
				_ = sdnotify.Stopping()
				stop()
				return
			}
		}
	}()

	// FIXME: setup application, start whatever services or listeners.
	//
	// NOTE: make sure when you are adding your application code, make sure none
	// of it blocks for an extended period of time, a common mistake is starting
	// a HTTP server on the main thread using [http.ListenAndServe] without
	// placing it in a go-routine.

	// Wait until the context is canceled, via a Interrupt signal.
	<-ctx.Done()

	// NOTE: you can add any application cleanup code here, or utilize `defer`
	// calls while starting your application.
}

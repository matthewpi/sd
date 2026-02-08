// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2025 Matthew Penner

//go:build linux

package sdnotify

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSdnotify(t *testing.T) {
	ctx := t.Context()

	// Override `getMonotonicUsec` to return a static value to make testing easier.
	getMonotonicUsec = func() int64 { return 4162392170 }

	// Clear the socket path just to be safe.
	socketPath = ""

	// Check that the address is as expected.
	if socketAddr() != nil {
		t.Errorf("expected socket address to be nil if `socketPath` is empty.")
	}

	// Create a new temporary path for us to setup a socket on.
	tmpDir, err := os.MkdirTemp(os.TempDir(), "nexavo")
	if err != nil {
		t.Fatal(fmt.Errorf("failed to create temporary directory: %w", err))
		return
	}
	socketPath = filepath.Join(tmpDir, "notify.sock")
	defer func() {
		_ = os.Remove(socketPath)
		_ = os.Remove(tmpDir)
	}()

	// Check that the address is as expected.
	addr := socketAddr()
	if expected, got := "unixgram", addr.Net; expected != got {
		t.Errorf("expected socket network to be \"%s\", but got \"%s\"", expected, got)
		return
	}
	if expected, got := socketPath, addr.Name; expected != got {
		t.Errorf("expected socket name to be \"%s\", but got \"%s\"", expected, got)
		return
	}

	// Start listening on the address.
	socket, err := net.ListenUnixgram(addr.Net, addr)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to start listening: %w", err))
		return
	}

	msg := make(chan []byte)
	go func(ctx context.Context, msg chan<- []byte) {
		defer socket.Close()
		context.AfterFunc(ctx, func() { _ = socket.SetDeadline(time.Now()) })

		buf := make([]byte, 16<<10)
		for {
			n, _, err := socket.ReadFromUnix(buf)
			if err != nil {
				if ctx.Err() != nil {
					// If there is an error from the context, exit without an error.
					return
				}
				t.Errorf("ReadFromUnix: %#v", err)
				continue
			}
			raw := buf[:n]

			t.Log(string(raw))
			msg <- raw
		}
	}(ctx, msg)

	for _, tc := range []struct {
		name   string
		fn     func() error
		expect []byte
	}{
		{
			name:   "Ready",
			fn:     Ready,
			expect: []byte(readyMessage),
		},
		{
			name:   "Reloading",
			fn:     Reloading,
			expect: []byte(reloadingMessage + "\n" + monotonicUsecPrefix + "4162392170"),
		},
		{
			name:   "STOPPING",
			fn:     Stopping,
			expect: []byte(stoppingMessage),
		},
	} {
		if err := tc.fn(); err != nil {
			t.Errorf("%s: %#v", tc.name, err)
			continue
		}

		if expected, got := tc.expect, <-msg; !bytes.Equal(expected, got) {
			t.Errorf("%s: expected \"%s\", but got \"%s\"", tc.name, expected, got)
			continue
		}
	}

	{
		data := []byte("Hello, world!")
		if err := StatusBytes(data); err != nil {
			t.Errorf("Status: %#v", err)
		} else if expected, got := prependString(statusPrefix, data), <-msg; !bytes.Equal(expected, got) {
			t.Errorf("Status: expected \"%s\", but got \"%s\"", expected, got)
		}
	}

	{
		testErr := errors.New("this is a test error\nwith a newline")
		// Notice how the new-line in the error gets replaced by a space, this is
		// necessary since `sd_notify` uses new-lines to separate key-value lines.
		expected := []byte(statusPrefix + "this is a test error with a newline\n" + errnoPrefix + "1")
		if err := Error(testErr, 1); err != nil {
			t.Errorf("Error: %#v", err)
		} else if got := <-msg; !bytes.Equal(expected, got) {
			t.Errorf("Error: expected \"%s\", but got \"%s\"", expected, got)
		}
	}
}

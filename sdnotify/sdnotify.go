// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2025 Matthew Penner

//go:build linux

package sdnotify

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
)

const (
	// readyMessage is the message to send when the application is ready.
	//
	// ref; https://www.freedesktop.org/software/systemd/man/latest/sd_notify.html#READY=1
	readyMessage = "READY=1"

	// reloadingMessage is the message to send when the application is reloading.
	//
	// ref; https://www.freedesktop.org/software/systemd/man/latest/sd_notify.html#RELOADING=1
	reloadingMessage = "RELOADING=1"

	// stoppingMessage is the message to send when the application is stopping.
	//
	// ref; https://www.freedesktop.org/software/systemd/man/latest/sd_notify.html#STOPPING=1
	stoppingMessage = "STOPPING=1"

	// statusPrefix is the prefix for notifying systemd of the application's
	// status. The argument for status is a freeform string and will be visible
	// in both the system's journal and via `systemctl status <NAME>.service`.
	//
	// ref; https://www.freedesktop.org/software/systemd/man/latest/sd_notify.html#STATUS=%E2%80%A6
	statusPrefix = "STATUS="

	// errnoPrefix is the prefix for sending an errno-style error code to
	// systemd when the application experiences an error.
	//
	// ref; https://www.freedesktop.org/software/systemd/man/latest/sd_notify.html#ERRNO=%E2%80%A6
	errnoPrefix = "ERRNO="

	// monotonicUsecPrefix is the prefix for sending a monotonic timestamp to
	// systemd.
	//
	// ref; https://www.freedesktop.org/software/systemd/man/latest/sd_notify.html#MONOTONIC_USEC=%E2%80%A6
	monotonicUsecPrefix = "MONOTONIC_USEC="
)

// socketPath is the path to the `sd_notify` socket. By default it will be set
// to the value of `os.Getenv("NOTIFY_SOCKET")`, but may be unset if necessary.
var socketPath = os.Getenv("NOTIFY_SOCKET")

// socketAddr returns the [*net.UnixAddr] for the `sd_notify` socket.
func socketAddr() *net.UnixAddr {
	if socketPath == "" {
		return nil
	}
	return &net.UnixAddr{
		Name: socketPath,
		Net:  "unixgram",
	}
}

// openSocket opens the `sd_notify` socket.
func openSocket() (*net.UnixConn, error) {
	addr := socketAddr()
	if addr == nil {
		return nil, nil
	}
	c, err := net.DialUnix(addr.Net, nil, addr)
	if err != nil {
		return nil, fmt.Errorf("sdnotify: unable to open NOTIFY_SOCKET: %w", err)
	}
	return c, nil
}

// sdnotify opens the `sd_notify` socket and sends the data in `payload` to it.
func sdnotify(payload []byte) error {
	c, err := openSocket()
	if c == nil || err != nil {
		return err
	}
	defer c.Close()
	if _, err = c.Write(payload); err != nil {
		return fmt.Errorf("sdnotify: failed to send message: %w", err)
	}
	return nil
}

// Notify sends data to the `sd_notify` socket.
//
// This can be used to send arbitrary messages to the `sd_notify` socket. Most
// applications should not use this and instead use the other functions provided
// by this package, such as [Ready], [Reloading], [Status], [Error], etc.
//
// If you are going to use this function directly, be careful. Do not chain
// multiple calls to [Notify] back-to-back, if you need to send multiple values,
// such as [Reloading] does (`RELOADING=1` and `MONOTONIC_USEC=...`), build a
// single byte-slice and call [Notify] once. Otherwise, systemd will treat each
// call to [Notify] as a separate message and issues may occur.
func Notify(payload []byte) error {
	return sdnotify(payload)
}

// Ready notifies `sd_notify` that the application is ready.
func Ready() error {
	return sdnotify([]byte(readyMessage))
}

// getMonotonicUsec holds a function that returns the current monotonic time,
// used to override the implementation during tests.
var getMonotonicUsec = nowMonotonic

// Reloading notifies `sd_notify` that the application is reloading.
//
// This function sends both `RELOADING=1` and `MONOTONIC_USEC=...` to systemd
// for proper support when running with `Type=notify-reload` instead of just
// `Type=notify`.
//
// This should be called right before the application starts reloading, once
// reloading is complete, [Ready] must be called unless an error occurs. If an
// error occurs during reloading, call [Error] instead of [Ready].
//
// Do your best to ensure that a failed reload doesn't break the application.
// It is better to error after a failed reload, but keep the application running
// with whatever config/settings were being used before the reload was triggered.
func Reloading() error {
	now, err := getMonotonicUsec()
	if err != nil {
		return fmt.Errorf("unable to get current monotonic time: %w", err)
	}
	usec := now.UnixMicro()

	var b bytes.Buffer
	b.WriteString(reloadingMessage)
	b.WriteByte('\n')
	b.WriteString(monotonicUsecPrefix)
	b.WriteString(strconv.FormatInt(usec, 10))
	return sdnotify(b.Bytes())
}

// Stopping notifies `sd_notify` that the application is stopping.
func Stopping() error {
	return sdnotify([]byte(stoppingMessage))
}

// Status sends a status message to `sd_notify`. The message will be visible in
// the both the system's journal and via `systemctl status <NAME>.service`.
func Status(msg string) error {
	return StatusBytes([]byte(msg))
}

// StatusBytes is like [Status] except that it takes a byte-slice instead of
// a string.
func StatusBytes(msg []byte) error {
	return sdnotify(prependString(statusPrefix, msg))
}

// Error sends an error message to `sd_notify`. The message will be visible in
// the system's journal and in `systemctl status <NAME>.service`.
//
// This function should only be used with an `errno`, without it, this function
// is the same as [Status].
func Error(err error, errno int) error {
	return ErrorBytes([]byte(err.Error()), errno)
}

// ErrorMessage is like [Error] except that it takes a string instead of
// an [error].
func ErrorMessage(msg string, errno int) error {
	return ErrorBytes([]byte(msg), errno)
}

// ErrorBytes is like [Error] except that it takes a byte-slice instead of
// an [error].
func ErrorBytes(msg []byte, errno int) error {
	var b bytes.Buffer
	b.WriteString(statusPrefix)
	b.Write(formatErrorMessage(msg))
	if errno > 0 {
		b.WriteByte('\n')
		b.WriteString(errnoPrefix)
		b.WriteString(strconv.Itoa(errno))
	}
	return sdnotify(b.Bytes())
}

// formatErrorMessage performs an efficient in-place replacement of new-lines
// with spaces instead of using [bytes.ReplaceAll] or [strings.ReplaceAll].
//
// This is used to avoid sending messages with new-lines to the `sd_notify`
// socket, which would not be interpreted correctly.
func formatErrorMessage(v []byte) []byte {
	for i, c := range v {
		if c == '\n' {
			v[i] = ' '
			continue
		}
	}
	return v
}

// prependString prepends a string (usually a constant) to a byte-slice.
func prependString(prefix string, data []byte) []byte {
	prefixLen := len(prefix)
	v := make([]byte, prefixLen+len(data))
	copy(v[:prefixLen], prefix)
	copy(v[prefixLen:], data)
	return v
}

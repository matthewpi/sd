// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2025 Matthew Penner

//go:build linux

package sdnotify

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	// watchdogMessage informs systemd to update the watchdog timestamp. This
	// message is used as a keep-alive ping when [WatchdogSec=] is configured
	// on the [systemd.service(5)] for this application.
	//
	// ref; https://www.freedesktop.org/software/systemd/man/latest/sd_notify.html#WATCHDOG=1
	//
	// [systemd.service(5)]: https://www.freedesktop.org/software/systemd/man/latest/systemd.service.html
	// [WatchdogSec=]: https://www.freedesktop.org/software/systemd/man/latest/systemd.service.html#WatchdogSec=
	watchdogMessage = "WATCHDOG=1"

	// watchdogTriggerMessage is used to inform systemd that an internal error
	// occurred.
	//
	// The result of calling this is the same as if [watchdogMessage] wasn't
	// sent to systemd in the proper interval.
	//
	// ref; https://www.freedesktop.org/software/systemd/man/latest/sd_notify.html#WATCHDOG=trigger
	watchdogTriggerMessage = "WATCHDOG=trigger"
)

// Watchdog informs systemd to update the watchdog timestamp. This is used as a
// keep-alive ping when [WatchdogSec=] is configured on the [systemd.service(5)]
// for this application.
//
// ref; https://www.freedesktop.org/software/systemd/man/latest/sd_notify.html#WATCHDOG=1
//
// [systemd.service(5)]: https://www.freedesktop.org/software/systemd/man/latest/systemd.service.html
// [WatchdogSec=]: https://www.freedesktop.org/software/systemd/man/latest/systemd.service.html#WatchdogSec=
func Watchdog() error {
	return sdnotify([]byte(watchdogMessage))
}

// WatchdogTrigger informs systemd that an internal error occurred.
//
// The result of calling this is the same as if [Watchdog] failed to send
// it's keep-alive in a given interval, except it will occur immediately,
// instead of after the interval was missed.
//
// ref; https://www.freedesktop.org/software/systemd/man/latest/sd_notify.html#WATCHDOG=trigger
func WatchdogTrigger() error {
	return sdnotify([]byte(watchdogTriggerMessage))
}

// WatchdogInterval returns the interval for the systemd watchdog if configured
// for the application.
//
// If the application is not running under systemd, or if the watchdog isn't
// configured, a duration of `0` and an error of `nil` will be returned.
//
// Applications wishing to implement support for systemd's watchdog, should
// create a [time.Ticker] (or similar) with the duration returned by this
// function, calling [Watchdog] at every tick.
func WatchdogInterval() (time.Duration, error) {
	// Get and parse `WATCHDOG_USEC` into a [time.Duration].
	wdUsec := os.Getenv("WATCHDOG_USEC")
	if wdUsec == "" {
		return 0, nil
	}
	usec, err := strconv.ParseInt(wdUsec, 10, 64)
	if err != nil {
		err = fmt.Errorf("sdnotify: unable to convert WATCHDOG_USEC to an integer: %w", err)
		return 0, nil
	}
	if usec < 1 {
		err = errors.New("sdnotify: WATCHDOG_USEC must be a positive integer")
		return 0, nil
	}
	// Convert the usec integer to a [time.Duration].
	d := time.Duration(usec) * time.Microsecond

	// Get and check `WATCHDOG_PID` against our PID.
	wdPid := os.Getenv("WATCHDOG_PID")
	if wdPid == "" {
		return 0, nil
	}
	pid, err := strconv.Atoi(wdPid)
	if err != nil {
		err = fmt.Errorf("sdnotify: unable to convert WATCHDOG_PID to an integer: %w", err)
		return 0, nil
	}
	if pid != os.Getpid() {
		return 0, nil
	}

	// `WATCHDOG_USEC` is set and `WATCHDOG_PID` matches our process id, so
	// return the duration and no error.
	return d, nil
}

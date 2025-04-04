// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2025 Matthew Penner

//go:build !linux

package sdnotify

import "time"

func Watchdog() error                          { return nil }
func WatchdogTrigger() error                   { return nil }
func WatchdogInterval() (time.Duration, error) { return 0, nil }

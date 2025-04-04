// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2025 Matthew Penner

//go:build !linux

package sdnotify

import "time"

// nowMonotonic returns the current time from the CLOCK_MONOTONIC clock.
func nowMonotonic() (t time.Time, err error) {
	return time.Now(), nil
}

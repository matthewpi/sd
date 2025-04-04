// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2025 Matthew Penner

package sdnotify

import (
	"testing"
	"time"
)

func TestNowMonotonic(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		now, err := nowMonotonic()
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(now)
	})

	t.Run("only moves forward", func(t *testing.T) {
		t1, err := nowMonotonic()
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(t1)

		time.Sleep(10 * time.Millisecond)

		t2, err := nowMonotonic()
		if err != nil {
			t.Error(err)
		}
		t.Log(t2)

		if !t1.Before(t2) {
			t.Error("t1 is after t2, this is not allowed to happen")
		}
	})
}

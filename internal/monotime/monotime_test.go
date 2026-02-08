// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2025 Matthew Penner

package monotime_test

import (
	"testing"
	"time"

	"github.com/matthewpi/sd/internal/monotime"
)

func TestNow(t *testing.T) {
	t.Run("back to back", func(t *testing.T) {
		t1 := monotime.Now()
		t2 := monotime.Now()
		if t1 > t2 {
			t.Error("t1 is after t2, this is not allowed to happen")
		}
	})

	t.Run("with sleep", func(t *testing.T) {
		t1 := monotime.Now()
		time.Sleep(10 * time.Millisecond)
		t2 := monotime.Now()
		if t1 > t2 {
			t.Error("t1 is after t2, this is not allowed to happen")
		}
	})
}

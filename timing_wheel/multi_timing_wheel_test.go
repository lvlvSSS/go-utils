package timing_wheel

import "testing"

func TestLevel(t *testing.T) {
	lvl, lvl_idx, idx := level(232)
	t.Logf("level : %d, idx in level : %d, idx: %d", lvl, lvl_idx, idx)

	lvl, lvl_idx, idx = level(256)
	t.Logf("level : %d, idx in level : %d, idx: %d", lvl, lvl_idx, idx)

	lvl, lvl_idx, idx = level(635)
	t.Logf("level : %d, idx in level : %d, idx: %d", lvl, lvl_idx, idx)

	lvl, lvl_idx, idx = level(1048790)
	t.Logf("level : %d, idx in level : %d, idx: %d", lvl, lvl_idx, idx)

	lvl, lvl_idx, idx = level(4194451)
	t.Logf("level : %d, idx in level : %d, idx: %d", lvl, lvl_idx, idx)
}

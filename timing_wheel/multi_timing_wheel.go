package timing_wheel

type jiffy int64

type TimingWheel struct {
	bucket_L1 [256]Node
	bucket_L2 [64]Node
	bucket_L3 [64]Node
	bucket_L4 [64]Node
	bucket_L5 [64]Node

	duration jiffy
	next     *TimingWheel
}

type Node struct {
	next, prev *Node
	do         func(callback func())
}

// level - like linux's timing wheel, l1 has 256 jiffies, l2~ln has 64 jiffies.
func level(duration uint64) (lvl, lvl_idx, idx int) {
	const L1 = 0xff
	const Ln = 0x3f
	idx = int(duration & L1)
	lvl = 0
	for cur := duration >> 8; cur&Ln > 0 || cur > Ln; cur = cur >> 6 {
		lvl_idx = int(cur&Ln) - 1
		lvl++
	}
	/*	cur := duration >> 8
		for {
			if cur&Ln <= 0 && cur < Ln {
				return
			}
			lvl++
			cur = cur >> 6
		}*/
	return
}

package timing_wheel

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var Unmodified = errors.New("state in timing wheel can't be modified now")

type state uint32

const (
	unstarted state = 1 << iota
	running
	finished
)

type TimingWheel struct {
	bucket *Bucket

	jiffy time.Duration
	next  *TimingWheel

	state state
	lock  sync.Mutex
}

func (wheel *TimingWheel) DefineJiffy(dur time.Duration) error {
	if state(atomic.LoadUint32((*uint32)(&wheel.state))) == unstarted {
		wheel.lock.Lock()
		defer wheel.lock.Unlock()
		if state(atomic.LoadUint32((*uint32)(&wheel.state))) == unstarted {
			wheel.jiffy = dur
		}
		return nil
	}
	return Unmodified
}

func (wheel *TimingWheel) Run() error {
	for {

		time.Sleep(wheel.jiffy / 4)
	}
	return nil
}

// location - like linux's timing wheel, l0 has 256 slots, l1~ln has 64 slots.
func location(duration uint64) (lvl, lvl_idx, idx int) {
	const L0 = 0xff
	const Ln = 0x3f
	idx = int(duration & L0)
	lvl = 0
	for cur := duration >> 8; cur&Ln > 0 || cur > Ln; cur = cur >> 6 {
		lvl_idx = int(cur & Ln)
		lvl++
	}
	/*	cur := jiffy >> 8
		for {
			if cur&Ln <= 0 && cur < Ln {
				return
			}
			lvl++
			cur = cur >> 6
		}*/
	return
}

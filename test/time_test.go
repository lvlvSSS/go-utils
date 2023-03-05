package test

import (
	"testing"
	"time"
)

func TestCurrentTimeStamp(t *testing.T) {
	duration, _ := time.ParseDuration("1650789964886ms")

	date := duration

	t.Logf("%v", date)

	t.Logf("%d", -1^(-1<<5))
	var i int64 = -31
	t.Logf("%b", i)
}

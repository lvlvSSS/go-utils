package timing_wheel

type Bucket struct {
	slots []Slot
	next  *Bucket
}

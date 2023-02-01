package timing_wheel

type Bucket struct {
	slots  []Slot
	number int // record the length of the slots.
	next   *Bucket
}

func (bucket *Bucket) Clear() {
	for idx, _ := range bucket.slots {
		bucket.slots[idx] = Slot{
			bucket: bucket,
		}
	}
}

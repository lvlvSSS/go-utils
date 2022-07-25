package skiplist

import (
	"math/rand"
)

/*
	Element is to describe the node of skiplist.
	single orient list.
*/
type Element struct {
	key []byte

	score float64

	data any

	list *SkipList
	// header's levels' count is equals to the count of nodes minus 1
	levels []*Element
}

func (element *Element) Key() []byte {
	return element.key
}

func (element *Element) Data() any {
	return element.data
}

func (element *Element) Score() float64 {
	return element.score
}

/*
	Level - return the count of levels.
*/
func (element *Element) Level() uint {
	if element.levels == nil {
		return 0
	} else {
		return uint(len(element.levels))
	}
}

func (element *Element) getScore() float64 {
	if element.key == nil {
		return -1
	}
	return calcScore(element.key)
}

func calcScore(key []byte) float64 {
	var hash uint64
	l := len(key)
	if l > 8 {
		l = 8
	}
	for i := 0; i < l; i++ {
		shift := uint(64 - 8 - i*8)
		hash |= uint64(key[i]) << shift
	}
	return float64(hash)
}

/*
	randLevel - get the random levels
*/
func randLevel() uint {
	i := 1
	for ; i <= MAX_LEVELS; i++ {
		//rand.Seed(int64(time.Now().Nanosecond()))
		if rand.Intn(2) == 0 {
			break
		}
	}
	return uint(i)
}

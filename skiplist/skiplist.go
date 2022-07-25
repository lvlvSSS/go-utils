package skiplist

import (
	"bytes"
	"math/rand"
	"time"
)

/*
	MAX_LEVELS - the max level of the skip list.
*/
const MAX_LEVELS = 32

type SkipList struct {
	header *Element
	len    uint ``
	level  uint
}

func (skipList *SkipList) Header() *Element {
	return skipList.header
}

func (skipList *SkipList) Level() uint {
	return skipList.level
}

func (skipList *SkipList) Len() uint {
	return skipList.len
}

/*
	New - Create a new skip list with header
*/
func New(depth uint) *SkipList {
	// make sure that the levels is positive and less than MAX_LEVELS.
	if depth <= 0 || depth > MAX_LEVELS {
		depth = MAX_LEVELS
	}

	skipList := SkipList{
		nil,
		0,
		0,
	}
	var header = skipList.createElement(depth, nil, nil)
	skipList.header = header
	// create the random seed for the randLevel()
	rand.Seed(int64(time.Now().Nanosecond()))
	return &skipList
}

/*
	Dispose - free all the children list in the skip list
*/
func (skipList *SkipList) Dispose() {
	if skipList.header == nil {
		return
	}
	current := skipList.header
	for index := skipList.level - 1; index >= 0; index-- {
		if current.levels[index] == nil {
			continue
		}
		// let each index of levels point to nil
		for next := current.levels[index]; next != nil; next = current.levels[index] {
			prev := current
			current = next
			prev.levels[index] = nil
		}
		current = skipList.header
	}
}

/*
	Search - search every index from top to bottom.
	return nil if not found.
*/
func (skipList *SkipList) Search(key []byte) *Element {
	prev := skipList.header
	index := skipList.level - 1
	for index >= 0 {
		for current := prev.levels[index]; current != nil; current = prev.levels[index] {
			if comp := compare(current, key, calcScore(key)); comp >= 0 {
				if comp == 0 {
					return current
				}
				// if the searched key is less than current, then the value maybe in lower index.
				break
			}
			prev = current
		}
		index--
	}
	return nil
}

/*
	Add - return true if insert a new element, return false if update the old element.
*/
func (skipList *SkipList) Add(key []byte, data any) bool {
	var successor = make([]*Element, MAX_LEVELS, MAX_LEVELS)
	prev := skipList.header

	// search the index that the key should insert to.
	for index := skipList.level - 1; index >= 0; index-- {
		successor[index] = prev
		for current := prev.levels[index]; current != nil; current = prev.levels[index] {
			if comp := compare(current, key, calcScore(key)); comp >= 0 {
				if comp == 0 {
					// return false if the key already exists.
					current.data = data
					return false
				}
				break
			}
			prev = current
			successor[index] = prev
		}
	}
	newLevel := randLevel()
	for i := newLevel; i < MAX_LEVELS; i++ {
		successor[i] = skipList.header
	}
	if newLevel > skipList.level {
		skipList.level = newLevel
	}
	skipList.len++

	newElement := skipList.createElement(newLevel, key, data)
	for index := newLevel - 1; index >= 0; index-- {
		newElement.levels[index] = successor[index].levels[index]
		successor[index] = newElement
	}
	return true
}

func (skipList *SkipList) Delete(key []byte) bool {
	// previous elements list.
	var successor = make([]*Element, MAX_LEVELS, MAX_LEVELS)
	prev := skipList.header
	var current *Element
	// search the previous element.
	for index := skipList.level - 1; index >= 0; index-- {
		successor[index] = prev
		for current = prev.levels[index]; current != nil; current = prev.levels[index] {
			if comp := compare(current, key, calcScore(key)); comp >= 0 {
				break
			}
			prev = current
			successor[index] = prev
		}
	}
	// return false if the key not found
	if current == nil || bytes.Compare(current.key, key) != 0 {
		return false
	}

	for index := current.Level() - 1; index >= 0; index-- {
		if successor[index].Level() < (index + 1) {
			continue
		}
		if successor[index].levels[index] == current {
			successor[index].levels[index] = current.levels[index]
			if skipList.header.levels[index] == nil {
				skipList.level--
			}
		}
	}
	skipList.len--
	return true
}

func (skipList *SkipList) createElement(levels uint, key []byte, data any) *Element {
	element := &Element{
		key:    key,
		data:   data,
		list:   skipList,
		levels: make([]*Element, levels, levels),
	}
	element.score = element.getScore()
	return element
}

func compare(element *Element, key []byte, score float64) int {
	if element.score == score {
		return bytes.Compare(element.key, key)
	}
	if (element.score - score) < 0 {
		return 1
	}
	return -1
}

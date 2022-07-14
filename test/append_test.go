package test

import (
	"fmt"
	"testing"
)

func TestAppend(t *testing.T) {
	arr := make([]int, 0, 5)
	ints := append(arr, 1)
	fmt.Printf("%v, %v, %v \n", len(ints), cap(ints), ints)
	print()
}

func print(args ...int) {
	fmt.Printf("args : %v, %v \n", len(args), args)
}

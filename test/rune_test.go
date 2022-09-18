package test

import "testing"

func TestRune(t *testing.T) {
	strRune := "中国"
	str := "china"

	t.Logf("strRune : %d , %d", len([]rune(strRune)), len(strRune))
	t.Logf("str : %d , %d", len([]rune(str)), len(str))
}

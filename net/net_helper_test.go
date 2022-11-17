package net

import "testing"

func TestIsIP(t *testing.T) {
	ip1 := "10.127.30.45"
	ip2 := "127.30.45"
	ip3 := "255.255.255.2567"
	t.Logf("%s: %v", ip1, IsIP(ip1))
	t.Logf("%s: %v", ip2, IsIP(ip2))
	t.Logf("%s: %v", ip3, IsIP(ip3))
}

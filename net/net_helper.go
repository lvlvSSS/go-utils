package net

import "regexp"

const IpRegex = `^([1-9]|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])(\.(\d|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])){3}$`

func IsIP(address string) bool {
	compile := regexp.MustCompile(IpRegex)
	return compile.MatchString(address)
}

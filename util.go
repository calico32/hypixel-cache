package main

import "strings"

func removeDashes(s string) string {
	return strings.Replace(s, "-", "", -1)
}

func insertUuidDashes(s string) string {
	s = removeDashes(s)

	if len(s) != 32 {
		panic("Invalid UUID")
	}

	return s[:8] + "-" + s[8:12] + "-" + s[12:16] + "-" + s[16:20] + "-" + s[20:]
}

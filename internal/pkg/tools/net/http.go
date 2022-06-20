package net

import "strings"

func GetReferDomain(requestURI string) string {
	tmpURI := requestURI
	if find := strings.Contains(requestURI, "//"); find {
		tmpURI = strings.Split(requestURI, "//")[1]
	}
	res := strings.Split(tmpURI, "/")[0]
	return res
}

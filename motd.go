package gomcstat

import (
	"regexp"
)

var colorCodeRegex = regexp.MustCompile(`(?i)[§&][0-9A-FK-OR]`)

func CleanMOTD(motd string) string {
	return colorCodeRegex.ReplaceAllString(motd, "")
}

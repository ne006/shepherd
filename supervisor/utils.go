package supervisor

import "strings"

func splitPath(spath string) []string {
	return strings.Split(spath, ".")
}

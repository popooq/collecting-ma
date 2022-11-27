package trimmer

import "strings"

func Trimmer(trimmed string) []string {
	trimmed = strings.TrimPrefix(trimmed, "/")
	trimmed = strings.TrimSuffix(trimmed, "/")
	splitted := strings.Split(trimmed, "/")
	return splitted
}

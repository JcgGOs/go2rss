package util

import "regexp"

var ReComment = regexp.MustCompile("<!--.*?-->")

func ClearComment(text string) string {
	return ReComment.ReplaceAllLiteralString(text, "")
}

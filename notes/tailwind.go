package notes

import (
	"strings"
)

var fullMatch = map[string]string{}

var prefixMatch = map[string]string{
	"h3": "text-lg font-semibold",
	"a":  "text-orange-700 hover:underline hover:decoration-dotted hover:text-orange-400",
	"th": "text-start",
	"td": "text-start",
}

func wrapForTailwind(content string) string {
	for k, v := range fullMatch {
		content = strings.ReplaceAll(content, tag(k), tagWithClass(k, v))
	}

	for k, v := range prefixMatch {
		content = strings.ReplaceAll(content, "<"+k, "<"+k+" class=\""+v+"\"")
	}

	return content
}

func tag(value string) string {
	return "<" + value + ">"
}

func tagWithClass(tag string, class string) string {
	return "<" + tag + " class=\"" + class + "\">"
}

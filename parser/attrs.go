package parser

import "strings"

func parseAttrs(raw string) (string, map[string]string) {
	attrs := make(map[string]string)

	parts := strings.Fields(raw)
	if len(parts) == 0 {
		return "", attrs
	}

	tag := parts[0]
	if idx := strings.IndexByte(tag, '='); idx != -1 {
		attrs["_default"] = tag[idx+1:]
		tag = tag[:idx]
	}

	for _, part := range parts[1:] {
		if idx := strings.IndexByte(part, '='); idx != -1 {
			attrs[part[:idx]] = part[idx+1:]
		}
	}

	return strings.ToLower(tag), attrs
}

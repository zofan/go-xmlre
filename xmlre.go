package xmlre

import (
	"regexp"
	"strings"
)

var (
	spaceRe = regexp.MustCompile(`\s+`)
	quoteRe = regexp.MustCompile(`["']`)
	tagRe   = regexp.MustCompile(`<([^\s>]+)(\s+[^>]+)?>`)
	attrRe  = regexp.MustCompile(`\s+([\w:-]+)\s*=`)
)

func Compile(pattern string) *regexp.Regexp {
	pattern = `(?is)` + pattern

	pattern = Format(pattern)

	pattern = strings.ReplaceAll(pattern, `\s*\s*`, `\s*`)
	pattern = strings.ReplaceAll(pattern, `\s+\s+`, `\s+`)

	return regexp.MustCompile(pattern)
}

func Format(pattern string) string {
	var parts []string

	var tmp string
	for _, c := range pattern {
		if c == '(' || c == ')' || c == '[' || c == ']' {
			if tmp != `` {
				parts = append(parts, tmp)
			}
			tmp = ``
		}

		tmp += string(c)
	}

	if tmp != `` {
		parts = append(parts, tmp)
	}

	for i, part := range parts {
		if i == 0 || parts[i-1][len(parts[i-1])-1] != '\\' {
			switch part[0] {
			case '(':
				parts[i] = formatGroup(part)
			case '[':
				parts[i] = formatCharset(part)
			case ')', ']':
				parts[i] = formatPart(part)
			}
		}
	}

	return strings.Join(parts, ``)
}

func formatGroup(pattern string) string {
	return formatPart(pattern)
}

func formatCharset(pattern string) string {
	pattern = quoteRe.ReplaceAllString(pattern, `"'`)

	pattern = strings.ReplaceAll(pattern, `0-9`, `\d`)
	pattern = strings.ReplaceAll(pattern, `А-я`, `А-Яа-яЁё`)

	return pattern
}

func formatPart(pattern string) string {
	matches := tagRe.FindAllStringSubmatch(pattern, -1)
	for _, m := range matches {
		pattern = strings.ReplaceAll(pattern, m[0], `\s*<`+m[1]+m[2]+`(?:\s+[^>]*)?>\s*`)
	}

	matches = attrRe.FindAllStringSubmatch(pattern, -1)
	for _, m := range matches {
		pattern = strings.ReplaceAll(pattern, m[0], `(?:[^>]+)? `+m[1]+`\s*=\s*`)
	}

	pattern = quoteRe.ReplaceAllString(pattern, `\s*["']\s*`)
	pattern = spaceRe.ReplaceAllString(pattern, `\s*`)

	return pattern
}

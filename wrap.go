package gg

import (
	"fmt"
	"strings"
	"unicode"
)

type measureStringer interface {
	MeasureString(s string) (w, h float64)
}

func splitOnSpace(x string) []string {
	var result []string
	pi := 0
	ps := false
	for i, c := range x {
		s := unicode.IsSpace(c)
		if s != ps && i > 0 {
			result = append(result, x[pi:i])
			pi = i
		}
		ps = s
	}
	result = append(result, x[pi:])
	return result
}

func wordWrap(m measureStringer, s string, width float64) []string {
	var result []string
	for _, line := range strings.Split(s, "\n") {
		fields := splitOnSpace(line)
		if len(fields)%2 == 1 {
			fields = append(fields, "")
		}
		x := ""
		for i := 0; i < len(fields); i += 2 {
			w, _ := m.MeasureString(x + fields[i])
			if w > width {
				if x == "" {
					at, st := truncateOverflow(m, fields[i], width)
					for st != "" {
						result = append(result, RemoveSpace(at))
						at, st = truncateOverflow(m, st, width)
					}
					result = append(result, RemoveSpace(at))
					x = ""
					continue
				} else {
					at, st := truncateOverflow(m, x, width)
					for st != "" {
						result = append(result, RemoveSpace(at))
						at, st = truncateOverflow(m, st, width)
					}
					result = append(result, RemoveSpace(at))
					x = ""
				}
			}
			x += fields[i] + fields[i+1]
		}
		if x != "" {
			result = append(result, RemoveSpace(x))
		}
	}
	for i, line := range result {
		result[i] = strings.TrimSpace(line)
	}
	return result
}

func truncateOverflow(m measureStringer, text string, width float64) (string, string) {
	t := []rune(text)
	for i := range t {
		w, _ := m.MeasureString(string(t[:i]))
		if w > width {
			fmt.Println(string(t[:i]), string(t[i:]))
			return string(t[:i]), string(t[i:])
		}
	}

	return text, ""
}

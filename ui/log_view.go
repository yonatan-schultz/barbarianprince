package ui

import (
	"strings"
)

const maxLogLines = 30

// wrapLine splits a string into lines of at most maxWidth runes each.
// Continuation lines preserve the leading indentation of the original line.
func wrapLine(s string, maxWidth int) []string {
	if maxWidth <= 0 {
		return []string{s}
	}
	runes := []rune(s)
	if len(runes) <= maxWidth {
		return []string{s}
	}

	// Measure leading indentation so continuation lines stay aligned.
	indent := 0
	for indent < len(runes) && runes[indent] == ' ' {
		indent++
	}
	// Only preserve indent when it leaves enough room for content.
	contIndent := []rune{}
	if indent > 0 && indent < maxWidth/2 {
		contIndent = []rune(strings.Repeat(" ", indent))
	}

	var result []string
	for len(runes) > maxWidth {
		// Try to break at a word boundary.
		cut := maxWidth
		for cut > maxWidth/2 && runes[cut] != ' ' {
			cut--
		}
		if cut <= maxWidth/2 {
			cut = maxWidth // no good break point, hard-cut
		}
		result = append(result, string(runes[:cut]))
		runes = runes[cut:]
		// Strip the break-space then re-apply indentation.
		for len(runes) > 0 && runes[0] == ' ' {
			runes = runes[1:]
		}
		if len(contIndent) > 0 && len(runes) > 0 {
			runes = append(contIndent, runes...)
		}
	}
	if len(runes) > 0 {
		result = append(result, string(runes))
	}
	return result
}

// RenderLog renders the narrative log panel.
// maxWidth is the visual width available per line; long messages are wrapped.
func RenderLog(log []string, height int, maxWidth int) string {
	if height <= 0 {
		height = 20
	}

	// Expand each log entry into wrapped lines, tagging each with its source index
	type taggedLine struct {
		text string
		idx  int // index into log slice (for age colouring)
	}
	var expanded []taggedLine
	for i, msg := range log {
		for _, l := range wrapLine(msg, maxWidth) {
			expanded = append(expanded, taggedLine{l, i})
		}
	}

	// Keep only the last `height` wrapped lines
	if len(expanded) > height {
		expanded = expanded[len(expanded)-height:]
	}

	// Determine the index of the most-recent original log entry visible
	newestIdx := -1
	if len(expanded) > 0 {
		newestIdx = expanded[len(expanded)-1].idx
	}

	var lines []string
	for _, tl := range expanded {
		age := newestIdx - tl.idx
		var styled string
		if age <= 3 {
			styled = StyleLog.Render(tl.text)
		} else {
			styled = StyleLogOld.Render(tl.text)
		}
		lines = append(lines, styled)
	}

	// Pad top with blank lines
	for len(lines) < height {
		lines = append([]string{""}, lines...)
	}

	return strings.Join(lines, "\n")
}
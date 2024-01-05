/*
Open Source Initiative OSI - The MIT License (MIT):Licensing

The MIT License (MIT)
Copyright Ralph Caraveo (deckarep@gmail.com)

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package pkg

import (
	"regexp"
	"strings"

	"github.com/deckarep/tips/pkg/ui"

	"github.com/charmbracelet/lipgloss"
)

var (
	httpPattern      = "(?P<http>HEAD|POST|GET|PUT|DELETE)"
	numberPattern    = `(?P<number>\b\d+(?:\.\d+)?\b)`
	pathPattern      = "(?P<filepath>~?/\\S+)"
	quotedStrPattern = "(?P<string>\"(?:.+?)?\")"
	tokenPattern     = "(?P<token>true|false|nil)"
	unixProcPattern  = "(?P<unixproc>\\w+\\[\\d+\\])"
	ipv4Pattern      = `(?P<ipv4>\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b)`
	ipv6Pattern      = `(?P<ipv6>(?:[a-f0-9:]+:+)+[a-f0-9]+)`

	// combinedPattern joins all regexes with the alternation | symbol.
	combinedPattern = strings.Join([]string{
		// WARNING: Order matters on how things get matched when there one pattern can be interpreted as a submatch
		// within a larger pattern.
		ipv6Pattern,
		ipv4Pattern,
		quotedStrPattern,
		tokenPattern,
		numberPattern,
		unixProcPattern,
		httpPattern,
		pathPattern,
	}, "|")

	superRegex = regexp.MustCompile(combinedPattern)

	// Capture and cache each capture group index by name.
	filepathIdx = superRegex.SubexpIndex("filepath")
	httpIdx     = superRegex.SubexpIndex("http")
	unixProcIdx = superRegex.SubexpIndex("unixproc")
	tokenIdx    = superRegex.SubexpIndex("token")
	numIdx      = superRegex.SubexpIndex("number")
	ipv4Idx     = superRegex.SubexpIndex("ipv4")
	ipv6Idx     = superRegex.SubexpIndex("ipv6")
)

// applyColorRules is responsible for colorizing matching segments in linear time after all regex sub-matches were
// identified. If no sub-matches were identified a default style is applied.
func applyColorRules(line string) string {
	matches := superRegex.FindAllStringSubmatchIndex(line, -1)

	var sb strings.Builder
	lastStartIdx := 0

	for _, m := range matches {
		// Segment before the match occurred.
		sb.WriteString(ui.Styles.Faint.Render(line[lastStartIdx:m[0]]))

		var style *lipgloss.Style

		// Add more token/styles here as order does not matter here.
		if m[tokenIdx*2] != -1 {
			style = &ui.Styles.Yellow
		} else if m[numIdx*2] != -1 {
			style = &ui.Styles.Cyan
		} else if m[httpIdx*2] != -1 {
			style = &ui.Styles.Magenta
		} else if m[unixProcIdx*2] != -1 {
			style = &ui.Styles.Red
		} else if m[filepathIdx*2] != -1 {
			style = &ui.Styles.Green
		} else if m[ipv4Idx*2] != -1 {
			style = &ui.Styles.Purple
		} else if m[ipv6Idx*2] != -1 {
			style = &ui.Styles.Purple
		}

		// Segment of the match itself, when we have a non-nil style.
		if style != nil {
			sb.WriteString(style.Render(line[m[0]:m[1]]))
			lastStartIdx = m[1]
		}
	}

	// Segment that is final to the end of the string.
	// NOTE: if no matches were found this will take care of the entire line.
	if lastStartIdx < len(line) {
		sb.WriteString(ui.Styles.Faint.Render(line[lastStartIdx:]))
	}

	return sb.String()
}

package app

import (
	"context"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/log"
	jsoniter "github.com/json-iterator/go"
	"io"
	"os"
	"regexp"
	"strings"
	"tips/pkg/ui"
)

var (
	// Colorization Rules regex (order matters, last one wins)
	colorRules = []regWithColor{
		// Keyword match: sudo|closed
		{reg: regexp.MustCompile(`\b(sudo|closed)\b`), color: ui.Styles.Red},
		// Any size int regex (no decimals, boundaries don't matter)
		{reg: regexp.MustCompile(`\b\d+\b`), color: ui.Styles.Blue},
		// IPV4 regex
		{reg: regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`), color: ui.Styles.Blue},
	}
)

type regWithColor struct {
	reg   *regexp.Regexp
	color lipgloss.Style
}

type segment struct {
	text      string
	colorized bool
}

func applyColorRule(rule regWithColor, segments []segment) []segment {
	var newSegments []segment

	for _, seg := range segments {
		if seg.colorized {
			newSegments = append(newSegments, seg)
			continue
		}

		locs := rule.reg.FindAllStringIndex(seg.text, -1)
		lastEnd := 0
		for _, loc := range locs {
			// Before the match
			if loc[0] > lastEnd {
				newSegments = append(newSegments, segment{text: seg.text[lastEnd:loc[0]]})
			}
			// The match itself
			newSegments = append(newSegments, segment{
				text:      rule.color.Render(seg.text[loc[0]:loc[1]]),
				colorized: true,
			})
			lastEnd = loc[1]
		}
		// After the last match
		if lastEnd < len(seg.text) {
			newSegments = append(newSegments, segment{text: seg.text[lastEnd:]})
		}
	}

	return newSegments
}

func applyColorRules(line string) string {
	segments := []segment{{text: line}}

	for _, rule := range colorRules {
		segments = applyColorRule(rule, segments)
	}

	// Reconstruct the line with default coloring for non-matching segments
	var sb strings.Builder
	for _, seg := range segments {
		if seg.colorized {
			sb.WriteString(seg.text)
		} else {
			sb.WriteString(ui.Styles.Faint.Render(seg.text))
		}
	}

	return sb.String()
}

func RenderLogLine(ctx context.Context, w io.Writer, idx int, hostname, line string) {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)

	if !cfg.NoColor {
		// Apply regex coloring/filtering.
		// Experiment: log syntax highlighter similar to https://github.com/bensadeh/tailspin
		// TODO: this logic still isn't quite right...but it's a start.
		line = applyColorRules(line)
	}

	hostPrefix := ui.Styles.Cyan.Render(fmt.Sprintf("%s (%d): ", hostname, idx))
	if _, err := fmt.Fprintln(w, hostPrefix+ui.Styles.Faint.Render(line)); err != nil {
		log.Error("error occurred during `Fprintln` to the local io.Writer:", err)
	}
}

func RenderIPs(ctx context.Context, tableView *GeneralTableView, w io.Writer) error {
	for _, devRow := range tableView.Rows {
		fmt.Println(devRow[2])
	}
	return nil
}

func RenderJson(ctx context.Context, tableView *GeneralTableView, w io.Writer) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	return json.NewEncoder(w).Encode(tableView)
}

func RenderTableView(ctx context.Context, tableView *GeneralTableView, w io.Writer) error {

	err := renderProlog(ctx, tableView, w)
	if err != nil {
		return err
	}

	err = renderBody(ctx, tableView, w)
	if err != nil {
		return err
	}

	err = renderEpilog(ctx, tableView, w)
	if err != nil {
		return err
	}

	return nil
}

func renderProlog(ctx context.Context, tableView *GeneralTableView, w io.Writer) error {
	// Show tailnet inspecting:
	fmt.Fprint(w, ui.Styles.Faint.Render("\nTailnet: "))
	fmt.Fprintln(w, ui.Styles.Bold.Render(tableView.Tailnet))

	fmt.Fprint(w, ui.Styles.Faint.Render("Query: "))
	fmt.Fprintln(w, ui.Styles.Bold.Render(tableView.Query))

	if len(tableView.DNSName) > 0 {
		fmt.Fprint(w, ui.Styles.Faint.Render(fmt.Sprintf("Self (%d): ", tableView.SelfView.Index)))
		fmt.Fprintln(w, ui.Styles.Bold.Render(tableView.SelfView.DNSName))
	}

	return nil
}

func renderBody(ctx context.Context, tableView *GeneralTableView, w io.Writer) error {
	const (
		SmHdrWidth = 6
	)

	var (
		re = lipgloss.NewRenderer(os.Stdout)

		CellStyle = re.NewStyle().Padding(0, 1).Width(20)

		HeaderStyle   = re.NewStyle().Foreground(ui.Colors.Purple).Bold(true).Align(lipgloss.Center)
		SmHeaderStyle = HeaderStyle.Copy().Width(SmHdrWidth).Align(lipgloss.Center)
		OddRowStyle   = CellStyle.Copy().Foreground(ui.Colors.Gray)
		EvenRowStyle  = CellStyle.Copy().Foreground(ui.Colors.LightGray)
	)

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == 0:
				if col == 0 {
					return SmHeaderStyle
				} else {
					return HeaderStyle
				}
			case row%2 == 0:
				if col == 0 {
					return EvenRowStyle.Copy().Width(SmHdrWidth)
				}
				return EvenRowStyle
			default:
				if col == 0 {
					return OddRowStyle.Copy().Width(SmHdrWidth)
				}
				return OddRowStyle
			}
		}).
		Headers(tableView.Headers...).
		Rows(tableView.Rows...)

	// Finally, render the table
	fmt.Fprintln(w, t)
	return nil
}

func renderEpilog(ctx context.Context, tableView *GeneralTableView, w io.Writer) error {
	// Render machine count and elapsed secs.
	fmt.Fprint(w, ui.Styles.Faint.Render("Total Machines: "))
	fmt.Fprint(w, ui.Styles.Bold.Render(fmt.Sprintf("%d", tableView.TotalMachines)))

	if tableView.APIElapsed.Seconds() > 1 {
		fmt.Fprint(w, ui.Styles.Faint.Render(", Elapsed (secs): "))
		fmt.Fprintln(w, ui.Styles.Bold.Render(fmt.Sprintf("%0.2f", tableView.APIElapsed.Seconds())))
	} else {
		fmt.Fprint(w, ui.Styles.Faint.Render(", Elapsed (ms): "))
		fmt.Fprintln(w, ui.Styles.Bold.Render(fmt.Sprintf("%0.2d", tableView.APIElapsed.Milliseconds())))
	}

	return nil
}

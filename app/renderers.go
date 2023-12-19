package app

import (
	"context"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/log"
	"io"
	"os"
	"regexp"
	"strings"
	"tips/pkg/ui"
)

var (
	// Colorization Rules regex (order matters, last one wins)
	colorRules = []*regexp.Regexp{
		// Keyword match
		regexp.MustCompile(`\bsudo\b`),
		// Any size int regex (no decimals, boundaries don't matter)
		regexp.MustCompile(`\b\d+\b`),
		// IPV4 regex
		regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`),
	}
)

func applyColorRule(rule *regexp.Regexp, line string) string {
	locs := rule.FindAllStringIndex(line, -1)

	if locs != nil {
		var sb strings.Builder
		lastEnd := 0

		for _, loc := range locs {
			// Append the part of the string before the match
			sb.WriteString(ui.Styles.Faint.Render(line[lastEnd:loc[0]]))

			// Append the match itself, colorized
			sb.WriteString(ui.Styles.Green.Render(line[loc[0]:loc[1]]))

			// Update the last processed index
			lastEnd = loc[1]
		}

		// Append the remainder of the string after the last match
		sb.WriteString(ui.Styles.Faint.Render(line[lastEnd:]))

		line = sb.String()
	}

	return line
}

func applyColorRules(line string) string {
	for _, rule := range colorRules {
		line = applyColorRule(rule, line)
	}
	return line
}

func RenderLogLine(ctx context.Context, w io.Writer, idx int, hostname, line string) {
	// Apply regex coloring/filtering.
	line = applyColorRules(line)

	// TODO: would be cool to add a log syntax highlighter like: https://github.com/bensadeh/tailspin
	hostPrefix := ui.Styles.Green.Render(fmt.Sprintf("%s (%d): ", hostname, idx))
	if _, err := fmt.Fprintln(w, hostPrefix+ui.Styles.Faint.Render(line)); err != nil {
		log.Error("error occurred during `Fprintln` to the local io.Writer:", err)
	}
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
	fmt.Fprintf(w, ui.Styles.Faint.Render("\nTailnet: "))
	fmt.Fprintln(w, ui.Styles.Bold.Render(tableView.Tailnet))

	fmt.Fprintf(w, ui.Styles.Faint.Render("Tailnet Query: "))
	fmt.Fprintln(w, ui.Styles.Bold.Render(tableView.Query))

	if len(tableView.DNSName) > 0 {
		fmt.Fprintf(w, ui.Styles.Faint.Render(fmt.Sprintf("Self (%d): ", tableView.SelfView.Index)))
		fmt.Fprintln(w, ui.Styles.Bold.Render(tableView.SelfView.DNSName))
	}

	return nil
}

func renderBody(ctx context.Context, tableView *GeneralTableView, w io.Writer) error {
	const (
		SmHdrWidth = 5
	)

	var (
		re = lipgloss.NewRenderer(os.Stdout)

		//Checkmark = re.NewStyle().SetString("✓").Foreground(green)
		CellStyle = re.NewStyle().Padding(0, 1).Width(20)

		// HeaderStyle is the lipgloss style used for the table headers.
		HeaderStyle   = re.NewStyle().Foreground(ui.Colors.Purple).Bold(true).Align(lipgloss.Center)
		SmHeaderStyle = HeaderStyle.Copy().Width(SmHdrWidth).Align(lipgloss.Center)
		// OddRowStyle is the lipgloss style used for odd-numbered table rows.
		OddRowStyle = CellStyle.Copy().Foreground(ui.Colors.Gray)
		// EvenRowStyle is the lipgloss style used for even-numbered table rows.
		EvenRowStyle = CellStyle.Copy().Foreground(ui.Colors.LightGray)
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

	fmt.Fprint(w, ui.Styles.Faint.Render(", Elapsed (secs): "))
	fmt.Fprintln(w, ui.Styles.Bold.Render(fmt.Sprintf("%0.2f", tableView.APIElapsed.Seconds())))

	return nil
}

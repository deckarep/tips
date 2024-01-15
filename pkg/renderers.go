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
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/deckarep/tips/pkg/ui"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/log"
	jsoniter "github.com/json-iterator/go"
)

func RenderRemoteSummary(ctx context.Context, w io.Writer, success, errors uint32, elapsed time.Duration) error {
	succStr := ui.Styles.Green.Render(fmt.Sprintf("%d", success))
	errStr := ui.Styles.Faint.Render(fmt.Sprintf("%d", errors))

	// Upon any errors lets highlight this fact.
	if errors > 0 {
		succStr = ui.Styles.Faint.Render(fmt.Sprintf("%d", success))
		errStr = ui.Styles.Red.Render(fmt.Sprintf("%d", errors))
	}

	summary := fmt.Sprintf("Finished: successes: %s, failures: %s, elapsed (secs): %.2f",
		succStr,
		errStr,
		elapsed.Seconds())

	if _, err := fmt.Fprintln(w, summary); err != nil {
		log.Error("error on `Fprintln` when writing elapsed time", "error", err)
	}
	return nil
}

func RenderLogLine(ctx context.Context, w io.Writer, idx int, isStdErr bool, hostname, alias, line string) {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)

	if len(alias) > 0 {
		hostname = alias
	}

	if !cfg.NoColor {
		// Apply regex coloring/filtering.
		// Experiment: log syntax highlighter similar to https://github.com/bensadeh/tailspin
		line = applyColorRules(line)
	}

	var descriptorSuffix string
	if isStdErr {
		descriptorSuffix = ui.Styles.Yellow.Render(fmt.Sprintf(" >%d", 2))
	} else {
		descriptorSuffix = ui.Styles.Cyan.Render(fmt.Sprintf(" >%d", 1))
	}

	hostPrefix := ui.Styles.Cyan.Render(fmt.Sprintf("%s%s (%d): ", hostname, descriptorSuffix, idx))
	if _, err := fmt.Fprintln(w, hostPrefix+ui.Styles.Faint.Render(line)); err != nil {
		log.Error("error occurred during `Fprintln` to the local io.Writer", "error", err)
	}
}

func RenderIPs(ctx context.Context, tableView *GeneralTableView, w io.Writer) error {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)
	ips := make([]string, 0, len(tableView.Rows))

	for _, devRow := range tableView.Rows {
		ips = append(ips, devRow[2])
	}

	if _, err := fmt.Fprintln(w, strings.Join(ips, cfg.IPsDelimiter)); err != nil {
		log.Error("error occurred during `Fprintln` to the local io.Writer", "error", err)
	}
	return nil
}

func RenderJson(ctx context.Context, tableView *GeneralTableView, w io.Writer) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	return json.NewEncoder(w).Encode(tableView)
}

func RenderASCIITableView(ctx context.Context, tableView *GeneralTableView, w io.Writer) error {
	// Create a new tabwriter.Writer. The 'minwidth', 'tabwidth', 'padding' and 'padchar' parameters can be adjusted to your needs.
	tw := tabwriter.NewWriter(w, 0, 4, 4, ' ', tabwriter.AlignRight)

	// Write the headers.
	for _, h := range tableView.Headers {
		fmt.Fprint(tw, h.Title+"\t") // Use \t as the column delimiter
	}
	fmt.Fprintln(tw) // End of the header line

	// Write all the rows.
	for _, row := range tableView.Rows {
		for _, col := range row {
			fmt.Fprint(tw, col+"\t")
		}
		fmt.Fprintln(tw) // End of the row
	}

	// Ensure all data is written
	return tw.Flush()
}

func RenderTableView(ctx context.Context, tableView *GeneralTableView, w io.Writer) error {

	err := renderTableProlog(ctx, tableView, w)
	if err != nil {
		return err
	}

	err = renderTableBody(ctx, tableView, w)
	if err != nil {
		return err
	}

	err = renderTableEpilog(ctx, tableView, w)
	if err != nil {
		return err
	}

	return nil
}

func renderTableProlog(ctx context.Context, tableView *GeneralTableView, w io.Writer) error {
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

func renderTableBody(ctx context.Context, tableView *GeneralTableView, w io.Writer) error {
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
		Headers(tableView.HeaderTitles()...).
		Rows(tableView.Rows...)

	// Finally, render the table
	if _, err := fmt.Fprintln(w, t); err != nil {
		log.Error("error occurred during `Fprintln` to the local io.Writer", "error", err)
	}
	return nil
}

func renderTableEpilog(ctx context.Context, tableView *GeneralTableView, w io.Writer) error {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)
	// Render machine count and elapsed secs.
	fmt.Fprint(w, ui.Styles.Faint.Render("Total Machines: "))
	fmt.Fprint(w, ui.Styles.Bold.Render(fmt.Sprintf("%d", tableView.TotalMachines)))

	var showElapsed = func(elapsed time.Duration) {
		if elapsed.Seconds() >= 1.0 {
			fmt.Fprint(w, ui.Styles.Faint.Render(", Elapsed: "))
			fmt.Fprintln(w, ui.Styles.Bold.Render(fmt.Sprintf("%0.2fsecs", elapsed.Seconds())))
		} else {
			fmt.Fprint(w, ui.Styles.Faint.Render(", Elapsed: "))
			fmt.Fprint(w, ui.Styles.Bold.Render(fmt.Sprintf("%0.2dms", elapsed.Milliseconds())))
		}
	}

	if cfg.TailscaleAPI.ElapsedTime > 0 {
		showElapsed(cfg.TailscaleAPI.ElapsedTime)
		fmt.Fprintln(w)
	} else if cfg.CachedElapsed > 0 {
		showElapsed(cfg.CachedElapsed)

		fmt.Fprint(w, ui.Styles.Faint.Render(", Cached: "))
		fmt.Fprintln(w, ui.Styles.Green.Render(ui.Symbols.Checkmark))
	}

	return nil
}

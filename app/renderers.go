package app

import (
	"context"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"io"
	"os"
	"tips/pkg/ui"
)

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
	var (
		re = lipgloss.NewRenderer(os.Stdout)

		//Checkmark = re.NewStyle().SetString("✓").Foreground(green)
		CellStyle = re.NewStyle().Padding(0, 1).Width(20)

		// HeaderStyle is the lipgloss style used for the table headers.
		HeaderStyle   = re.NewStyle().Foreground(ui.Colors.Purple).Bold(true).Align(lipgloss.Center)
		SmHeaderStyle = HeaderStyle.Copy().Width(4).Align(lipgloss.Center)
		// OddRowStyle is the lipgloss style used for odd-numbered table rows.
		OddRowStyle = CellStyle.Copy().Foreground(ui.Colors.Gray)
		// EvenRowStyle is the lipgloss style used for even-numbered table rows.
		EvenRowStyle = CellStyle.Copy().Foreground(ui.Colors.LightGray)
	)

	const (
		SmHdrWidth = 4
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

// TODO: This stuff shouldn't be here...all of this should happen in a Processor
//func getRowField(enrichedResults map[string]tailscale_cli.DeviceInfo, idx int, d tailscale.Device) []string {
//	// You can also add tables row-by-row
//
//	var (
//		version = fmt.Sprintf("%s - %s", strings.Split(d.ClientVersion, "-")[0], d.OS)
//
//		timeAgo = humanize.Time(d.LastSeen.Time)
//		// Remove all tag: prefixes, and join the tags as a comma delimited string.
//		tags = strings.Replace(strings.Join(d.Tags, ","), "tag:", "", -1)
//	)
//
//	seenAgo := timeAgo
//	//if strings.Contains(seenAgo, "seconds") {
//	//	// https://github.com/tailscale/tailscale/pull/3534/files
//	//	seenAgo = fmt.Sprintf("• %s", seenAgo)
//	//}
//
//	if len(enrichedResults) > 0 {
//		if enrichedDev, ok := enrichedResults[d.NodeKey]; ok && enrichedDev.Online {
//			seenAgo = fmt.Sprintf("%s now", ui.Styles.Green.Render(ui.Dot))
//		}
//		return []string{strconv.Itoa(idx), d.Hostname, d.Addresses[0], tags, d.User, version, seenAgo}
//	}
//	return []string{strconv.Itoa(idx), d.Hostname, d.Addresses[0], tags, d.User, version, seenAgo}
//}
//
//func getHeaders(enrichedResults map[string]tailscale_cli.DeviceInfo) []string {
//	if len(enrichedResults) > 0 {
//		return []string{"No", "Machine", "Address", "Tags", "User", "Version", "LastSeen"}
//	}
//	return []string{"No", "Machine", "Address", "Tags", "User", "Version", "LastSeen"}
//}

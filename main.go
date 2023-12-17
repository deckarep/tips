/*
Open Source Initiative OSI - The MIT License (MIT):Licensing

The MIT License (MIT)
Copyright (c) 2023 - 2024 Ralph Caraveo (deckarep@gmail.com)

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

package main

import (
	"github.com/charmbracelet/log"
	"github.com/tailscale/tailscale-client-go/tailscale"
	"os"
	"time"
	"tips/cmd"
	"tips/pkg"
)

// Alternatively shell out to this command to get the Online status flag
// /Applications/Tailscale.app/Contents/MacOS/Tailscale status --json | jq .
// Totally janky.
// curl -s -u "$tips_api_key:" -XGET "https://api.tailscale.com/api/v2/tailnet/deckarep@gmail.com/devices" | jq .

/*
TODO: cache feature --nocache flag to force a refresh
TODO: show a nice useragent format: MyApp/1.0.0 (Windows NT 10.0; Win64; x64)
TODO: -ipv6 flag for ipv6 results
TODO: output everything as JSON
TODO: nail down the default table header/columns, provide config to enable disable for a user.
TODO: themify for TailScale, or enable theming feature
TODO: allow for easy filtering
TODO: ssh into node
TODO: similar to ssh, what about curl/http requests to all selected nodes?
TODO: default client timeout, configurable
TODO: support --ssh (single) or --csshx all hosts!
	- brew install parera10/csshx/csshx (original csshx is broken on recent macos)
TODO: allow for tailing logs of one or more boxes, ensure tailing logs works like knife
TODO: show online status (like the admin console)
TODO: show the count of all nodes returned near the bottom.
TODO: filter glob syntax: tips * 'hostname', tips blade* 'hostname', tips tag:peanuts 'hostname'
TODO: default sorting (hostname, then address)
TODO: slice syntax:
	- tips blade[5:10] 'hostname' - returns hosts 5-10 (can't remember if it should be inclusive or not
	- tips blade[5:] 'hostname' - returns host 5 on up
TODO: concurrency flag -c with some sane default like 5
TODO: enrich the API output with more details like "Online" status flag from the cli, when we're on a tailscale network.
*/

/*
BUG: Why is LastSeen not as recent as what's on the admin page? Perhaps heavy caching for API access.
	- maybe use the tsnet application for more realtime knowledge of what's going on: https://tailscale.com/kb/1244/tsnet/
*/

const (
	defaultClientTimeout = time.Second * 5
)

var (
	tailnet           = "deckarep@gmail.com"
	oauthClientID     = os.Getenv("tips_client_id")
	oauthClientSecret = os.Getenv("tips_secret")

	api_key = os.Getenv("tips_api_key")

	client = func() *tailscale.Client {
		client, err := tailscale.NewClient(
			api_key, // When doing oauth, this field must be blank!!!
			tailnet,
			//tailscale.WithOAuthClientCredentials(oauthClientID, oauthClientSecret, nil),
			tailscale.WithUserAgent(pkg.UserAgent),
		)
		if err != nil {
			log.Fatal("failed to create client with err: ", err)
		}
		return client
	}()
)

func main() {
	cmd.Execute()

	//if true {
	//	return
	//}
	//
	//viper.AddConfigPath(".")
	//viper.SetConfigName("tips.config") // Register config file name (no extension)
	//viper.SetConfigType("toml")        // Look for specific type
	//if err := viper.ReadInConfig(); err != nil {
	//	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
	//		log.Warn("config file was not found, falling back to defaults.")
	//	} else {
	//		log.Fatal("issue reading viper config file exiting with err: ", err)
	//	}
	//}
	//
	//cacheDuration := viper.GetDuration("settings.cli.cache_duration")
	//fmt.Println(cacheDuration)
	//
	//enrichedResults, err := tailscale_cli.GetDevicesStatuses()
	//if err != nil {
	//	fmt.Println("failed to get results: ", err)
	//}
	//
	//startTime := time.Now()
	//
	//// Show tailnet inspecting:
	//fmt.Printf(ui.Styles.Faint.Render("\nTailnet: "))
	//fmt.Println(ui.Styles.Bold.Render(tailnet))
	//
	//// Show user's query:
	////fmt.Println("•••\n•••\n•••") // create tailscale logo???
	//fmt.Printf(ui.Styles.Faint.Render("Tailnet Query: "))
	//fmt.Println(ui.Styles.Bold.Render("tips blade*"))
	//
	//if len(enrichedResults) > 0 {
	//	fmt.Printf(ui.Styles.Faint.Render(fmt.Sprintf("Self (%d): ", 2)))
	//	fmt.Println(ui.Styles.Bold.Render(pkg.SelfDevice(enrichedResults).DNSName))
	//}
	//
	//ctx := context.Background()
	//ctxTimeOut, cancelTimeout := context.WithTimeout(ctx, defaultClientTimeout)
	//defer cancelTimeout()
	//devList, err := client.Devices(ctxTimeOut)
	//if err != nil {
	//	log.Fatal("failed to fetch devices with err: ", err)
	//}
	//
	//// Sort device list via generic SortFunc.
	//slices.SortFunc(devList, func(a, b tailscale.Device) int {
	//	// Needs additional work.
	//	if n := cmp.Compare(a.Name, b.Name); n != 0 {
	//		return n
	//	}
	//
	//	return cmp.Compare(a.Name, b.Name)
	//})
	//
	//rows := [][]string{}
	//
	//var (
	//	re = lipgloss.NewRenderer(os.Stdout)
	//
	//	//Checkmark = re.NewStyle().SetString("✓").Foreground(green)
	//	CellStyle = re.NewStyle().Padding(0, 1).Width(20)
	//
	//	// HeaderStyle is the lipgloss style used for the table headers.
	//	HeaderStyle   = re.NewStyle().Foreground(ui.Colors.Purple).Bold(true).Align(lipgloss.Center)
	//	SmHeaderStyle = HeaderStyle.Copy().Width(4).Align(lipgloss.Center)
	//	// OddRowStyle is the lipgloss style used for odd-numbered table rows.
	//	OddRowStyle = CellStyle.Copy().Foreground(ui.Colors.Gray)
	//	// EvenRowStyle is the lipgloss style used for even-numbered table rows.
	//	EvenRowStyle = CellStyle.Copy().Foreground(ui.Colors.LightGray)
	//)
	//
	//const (
	//	SmHdrWidth = 4
	//)
	//
	//t := table.New().
	//	Border(lipgloss.NormalBorder()).
	//	BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
	//	StyleFunc(func(row, col int) lipgloss.Style {
	//		switch {
	//		case row == 0:
	//			if col == 0 {
	//				return SmHeaderStyle
	//			} else {
	//				return HeaderStyle
	//			}
	//		case row%2 == 0:
	//			if col == 0 {
	//				return EvenRowStyle.Copy().Width(SmHdrWidth)
	//			}
	//			return EvenRowStyle
	//		default:
	//			if col == 0 {
	//				return OddRowStyle.Copy().Width(SmHdrWidth)
	//			}
	//			return OddRowStyle
	//		}
	//	}).
	//	Headers(getHeaders(enrichedResults)...).
	//	Rows(rows...)
	//
	//for idx, d := range devList {
	//
	//	t.Row(getRowField(enrichedResults, idx, d)...)
	//
	//	//fmt.Println(d.NodeKey, d.MachineKey)
	//}
	//
	//// Render table.
	//fmt.Println(t)
	//
	//// Render machine count and elapsed secs.
	//fmt.Printf(ui.Styles.Faint.Render("Total Machines: "))
	//fmt.Printf(ui.Styles.Bold.Render(fmt.Sprintf("%d", len(devList))))
	//
	//fmt.Printf(ui.Styles.Faint.Render(", Elapsed (secs): "))
	//elapsedTime := time.Since(startTime)
	//fmt.Println(ui.Styles.Bold.Render(fmt.Sprintf("%.2f", elapsedTime.Seconds())))
}

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

//func getHeaders(enrichedResults map[string]tailscale_cli.DeviceInfo) []string {
//	if len(enrichedResults) > 0 {
//		return []string{"No", "Machine", "Address", "Tags", "User", "Version", "LastSeen"}
//	}
//	return []string{"No", "Machine", "Address", "Tags", "User", "Version", "LastSeen"}
//}

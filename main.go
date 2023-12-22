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
	"tips/cmd"
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
	- Not based filter: tag:!peanuts
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

func main() {
	cmd.Execute()
}

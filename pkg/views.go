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

import "time"

// TailnetView has everything known about a Tailnet
type TailnetView struct {
	Tailnet       string
	TotalMachines int
}

type ContextView struct {
	Query      string
	APIElapsed time.Duration
	CLIElapsed time.Duration
}

type SelfView struct {
	Index   int
	DNSName string
}

// DevicesView has everything needed to be rendered.
type DevicesView struct {
}

type DevicesTable struct {
	TailnetView
	Devices *DevicesView
}

type GeneralTableView struct {
	ContextView
	TailnetView
	SelfView
	Headers []Header
	Rows    [][]string
}

func (g *GeneralTableView) HeaderTitles() []string {
	var names []string
	for _, hdr := range g.Headers {
		names = append(names, hdr.Title)
	}
	return names
}

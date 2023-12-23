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
	Headers []string
	Rows    [][]string
}

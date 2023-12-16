package app

// TailnetView has everything known about a Tailnet
type TailnetView struct {
	Tailnet string
}

// DevicesView has everything needed to be rendered.
type DevicesView struct {
}

type DevicesTable struct {
	TailnetView
	Devices *DevicesView
}

type GeneralTableView struct {
	TailnetView
	Headers []string
	Rows    [][]string
}

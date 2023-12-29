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

package ui

import "github.com/charmbracelet/lipgloss"

type symbols struct {
	Checkmark string
	Dot       string
}

type styleTypes struct {
	Bold  lipgloss.Style
	Faint lipgloss.Style

	Black   lipgloss.Style
	Blue    lipgloss.Style
	Cyan    lipgloss.Style
	Green   lipgloss.Style
	Magenta lipgloss.Style
	Red     lipgloss.Style
	White   lipgloss.Style
	Yellow  lipgloss.Style
}

type colorTypes struct {
	Black     lipgloss.Color
	Blue      lipgloss.Color
	Cyan      lipgloss.Color
	Gray      lipgloss.Color
	Green     lipgloss.Color
	LightGray lipgloss.Color
	Magenta   lipgloss.Color
	Purple    lipgloss.Color
	Red       lipgloss.Color
	White     lipgloss.Color
	Yellow    lipgloss.Color
}

var (
	Symbols = symbols{
		Checkmark: "✔",
		Dot:       "•",
	}

	Colors = colorTypes{
		Black:     "16",
		Blue:      "12",
		Cyan:      "14",
		Gray:      "245",
		Green:     "10", // "32", "99",
		LightGray: "241",
		Magenta:   "5",
		Purple:    "99",
		Red:       "9",
		White:     "15",
		Yellow:    "11",
	}

	Styles = styleTypes{
		Bold: lipgloss.NewStyle().Bold(true),
		//Foreground(lipgloss.Color("#FAFAFA")).
		//Background(lipgloss.Color("#7D56F4")).
		//PaddingTop(2).
		//PaddingLeft(4).
		//Width(22)

		Faint: lipgloss.NewStyle().Faint(true),
		Green: lipgloss.NewStyle().
			Bold(true).
			Foreground(Colors.Green),
		Blue: lipgloss.NewStyle().
			Bold(true).
			Foreground(Colors.Blue),
		Cyan: lipgloss.NewStyle().
			Bold(true).
			Foreground(Colors.Cyan),
		Magenta: lipgloss.NewStyle().
			Bold(true).
			Foreground(Colors.Magenta),
		White: lipgloss.NewStyle().
			Bold(true).
			Foreground(Colors.Yellow),
		Yellow: lipgloss.NewStyle().
			Bold(true).
			Foreground(Colors.Yellow),
		Red: lipgloss.NewStyle().
			Bold(true).
			Foreground(Colors.Red),
		Black: lipgloss.NewStyle().
			Bold(true).
			Foreground(Colors.Black),
	}
)

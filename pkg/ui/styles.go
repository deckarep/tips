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

const (
	Checkmark = "✔"
	Dot       = "•"
)

type styleTypes struct {
	Bold  lipgloss.Style
	Faint lipgloss.Style
	Green lipgloss.Style
}

type colorTypes struct {
	Purple    lipgloss.Color
	Gray      lipgloss.Color
	LightGray lipgloss.Color
	Green     lipgloss.Color
}

var (
	Colors colorTypes = colorTypes{
		Purple:    lipgloss.Color("99"),
		Gray:      lipgloss.Color("245"),
		LightGray: lipgloss.Color("241"),
		Green:     lipgloss.Color("99"),
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
			Foreground(lipgloss.Color("#00ff00")),
	}
)

var ()

package main

import (
	"github.com/charmbracelet/bubbles/table"
	"strconv"
)

func (m model) GetTable() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func (m *model) ActualizeTable() {

	rows := []table.Row{}

	for _, timeLog := range m.tabLogs[m.activeTab] {
		str := table.Row{
			strconv.Itoa(timeLog.day),
			strconv.Itoa(timeLog.month),
			strconv.Itoa(timeLog.year),
			strconv.Itoa(timeLog.week),
		}
		rows = append(rows, str)
	}

	columns := []table.Column{
		{Title: "Day", Width: 6},
		{Title: "Month", Width: 6},
		{Title: "Year", Width: 6},
		{Title: "Week", Width: 6},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		// table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(tableHeaderBorderStyle).
		BorderForeground(tableHeaderBorderForeground).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(tableSelectedForeground).
		Background(tableSelectedBackground).
		Bold(false)
	t.SetStyles(s)

	m.table = t
}
